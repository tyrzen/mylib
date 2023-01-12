package main

import (
	"fmt"

	"github.com/delveper/mylib/lib/env"
	"github.com/delveper/mylib/lib/tokay"
)

func main() {
	env.Load()
	err := tokay.Parse("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWiOiJtZSIsImV4cCI6MTY3MzU2MjkxOSwiaWF0IjoxNjczNTYxMTE5LCJpc3MiOiJteWxpYiJ9.eFUe-noKO_Xy7ylKZcp6mL2lP56B1qBb-jmsOecKGnE")
	fmt.Println(err)
}
