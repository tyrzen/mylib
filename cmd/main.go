package main

import (
	"log"
	"os"

	"mylib/app/ent"
	"mylib/cfg"
	"mylib/lgr"
)

func main() {
	if err := cfg.Load(); err != nil {
		log.Fatalf("failed loading environment variables: %v", err)
	}

	logLvl := os.Getenv("LOG_LEVEL")

	var log ent.Logger = lgr.New(logLvl)

	log.Infof("Logger set up with log level %s", logLvl)

}
