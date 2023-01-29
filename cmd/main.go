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
	run()
	/*	env.LoadVars()

		alg := os.Getenv("JWT_ALG")
		key := os.Getenv("JWT_KEY")

		exp, err := time.ParseDuration(os.Getenv("JWT_ACCESS_EXP"))
		if err != nil {
			log.Fatalf("error parsing refresh token expirity: %v", err)
		}

		id := xid.New().String()

		payload := models.AccessToken{ID: id, RefreshTokenID: "<blank>", Expiry: exp}

		token, err := tokay.Make[models.AccessToken](alg, key, exp, payload)
		log.Println(token)
		if err != nil {
			log.Fatalf("error making token: %v", err)
		}

		data, err := tokay.Parse[models.AccessToken](token, key)
		if err != nil {
			log.Fatalf("error parsing token: %v", err)
		}

		log.Printf("%T\n", data)
		log.Printf("%+v", data)

	*/
}

func run() {
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

	readerLogic := usecases.NewReader(readerRepo, tokenRepo)

	readerREST := rest.NewReader(readerLogic, logger)

	router := rest.NewRouter(readerREST.Route)

	logger.Infof("Router successfully created.")

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
