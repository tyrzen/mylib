package main

import (
	"log"
	"os"

	"github.com/delveper/mylib/app/models"
	"github.com/delveper/mylib/app/presenters/rest"
	repo "github.com/delveper/mylib/app/repository/psql"
	"github.com/delveper/mylib/app/repository/rds"
	"github.com/delveper/mylib/app/usecases"
	"github.com/delveper/mylib/lib/banderlog"
	"github.com/delveper/mylib/lib/env"
	"github.com/delveper/mylib/mig"
)

func main() {
	Run()
}

func Run() {
	if err := env.LoadVars(); err != nil {
		log.Printf("Failed to load environment variables: %+v", err)
		return
	}

	var logger models.Logger = banderlog.New()
	defer func() {
		if err := logger.Flush(); err != nil {
			log.Printf("Failed flush logger: %+v", err)
		}
	}()
	logger.Infof("Logger set up with level: %s", logger.Level())

	conn, err := repo.Connect()
	if err != nil {
		logger.Errorf("Failed connecting to repo: %+v", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Warnf("Failed closing repo connection: %+v", err)
		}
	}()
	logger.Infof("Connection to repo established.")

	migration := mig.New()
	migration.SetLogger(logger)

	if err := migration.Run(conn); err != nil {
		logger.Errorf("Failed making migrations: %+v", err)
		return
	}

	client, err := rds.Connect()
	if err != nil {
		logger.Errorf("Failed connecting to session repo: %+v", err)
		return
	}

	defer func() {
		if err := client.Close(); err != nil {
			logger.Warnf("Failed closing session repo connection: %+v", err)
		}
	}()

	readerRepo := repo.NewReader(conn)
	bookRepo := repo.NewBook(conn)
	tokenRepo := rds.NewToken(client)
	authorRepo := repo.NewAuthor(conn)

	logger.Infof("Repository layer initialized.")

	readerLogic := usecases.NewReader(readerRepo, tokenRepo)
	bookLogic := usecases.NewBook(bookRepo, authorRepo)

	logger.Infof("Usecase layer initialized.")

	readerREST := rest.NewReader(readerLogic, logger)
	bookREST := rest.NewBook(bookLogic, logger)

	logger.Infof("RESTish layer initialized.")

	router := rest.NewRouter(
		readerREST.Route,
		bookREST.Route,
	)

	logger.Infof("Routes created successfully.")

	handler := rest.ChainMiddlewares(router,
		rest.WithRequestID,
		rest.WithLogRequest(logger),
	)

	logger.Infof("Pre-middleware set up.")

	srv, err := rest.NewServer(handler)
	if err != nil {
		logger.Errorf("Failed initializing server: %+v", err)
		return
	}

	if err := srv.Run(); err != nil {
		logger.Errorf("Failed running server: %+v", err)
	}

	logger.Infof("Server started on the port: %s", os.Getenv("SRV_PORT"))
}
