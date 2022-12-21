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

	log := ent.Logger(lgr.New(os.Getenv("LOG_LEVEL")))

	_ = log
}
