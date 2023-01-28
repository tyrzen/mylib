package main

import (
	"log"
	"os"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/logic"
	repo "github.com/delveper/mylib/app/repo/psql"
	"github.com/delveper/mylib/app/repo/rds"
	"github.com/delveper/mylib/app/trans/rest"
	"github.com/delveper/mylib/lib/banderlog"
	"github.com/delveper/mylib/lib/env"
	"github.com/delveper/mylib/lib/jwt"
	"github.com/delveper/mylib/mig"
)

func main() {
	run()
	/*	env.LoadVars()
		val := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c`
		logic := jwt.New()
		key := os.Getenv("JWT_KEY")
		data, err := logic.Parse(val, key)
		log.Println(err)
		log.Println(key)
		log.Printf("%T\n", data)
		log.Printf("%+v", data)
	*/
}

func run() {
	if err := env.LoadVars(); err != nil {
		log.Printf("Failed to load environment variables: %+v", err)
		return
	}

	var logger ent.Logger = banderlog.New()
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
	tokenRepo := rds.NewToken(client)
	jwtLogic := jwt.New()

	readerLogic := logic.NewReader(readerRepo, tokenRepo, jwtLogic)

	readerREST := rest.NewReader(readerLogic, logger)

	router := rest.NewRouter(readerREST.Route)

	logger.Infof("Router successfully created.")

	handler := rest.ChainMiddlewares(router, rest.WithRequestID, rest.WithLogRequest(logger))

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
