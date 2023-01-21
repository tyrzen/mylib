package main

import (
	"time"

	"github.com/delveper/mylib/app/repo/rds"
	"github.com/delveper/mylib/lib/env"
)

func main() {
	env.LoadVars()
	red := rds.Conn()

	ctx := context.Background()

	if err := red.Set(ctx, "bla bla", "1", 10*time.Minute).Err(); err != nil {
		log.Println(err)
	}
}
