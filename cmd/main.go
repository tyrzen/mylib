package main

import (
	"log"

	"mylib/app/ent"
	"mylib/cfg"
	"mylib/lgr"
)

func main() {
	if err := cfg.Load(); err != nil {
		log.Fatalf("Failed to load environment variables: %+v", err)
	}

	var logger ent.Logger = lgr.New()
	defer func() {
		if err := logger.Flush(); err != nil {
			log.Printf("Failed flush logger: %+v", err)
		}
	}()

	logger.Infof("Logger set up with level: %s", logger.Level())
}
