package main

import (
	"log"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/lgr"
	repo "github.com/delveper/mylib/app/repo/psql"
	"github.com/delveper/mylib/app/trans/rest"
	"github.com/delveper/mylib/env"
	"github.com/delveper/mylib/mig"
)

func main() {
	if err := env.Load(); err != nil {
		log.Printf("Failed to load environment variables: %+v", err)
		return
	}

	var logger ent.Logger = lgr.New()

	defer func() {
		if err := logger.Flush(); err != nil {
			log.Printf("Failed flush logger: %+v", err)
		}
	}()
	logger.Infof("Logger set up with level: %s", logger.Level())

	conn, err := repo.Connect()
	if err != nil {
		logger.Errorf("Failed connecting to repo: %+v", err)
		return
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Warnf("Failed closing repo connection: %+v", err)
		}
	}()
	logger.Infof("Connection to repo was established.")

	migration := mig.New()
	migration.SetLogger(logger)

	if err := migration.Run(conn); err != nil {
		logger.Errorf("Failed make migrations: %+v", err)
		return
	}

	reader := rest.NewReader(nil, logger)

	rtr := rest.NewRouter(reader.Route)

	srv, err := rest.NewServer(rtr)
	if err != nil {
		logger.Errorf("Failed initializing server: %+v", err)
		return
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Errorf("Failed running server: %+v", err)
	}
}
