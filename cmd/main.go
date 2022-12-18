package main

import (
	"fmt"

	"mylib/lib/revalid" // custom validator based on RegExp
)

type User struct {
	Name     string `revalid:".+"`
	Password string `revalid:".{8,}"`
}

func main() {
	usr := User{
		Name:     "Joe",
		Password: "123452134",
	}
	err := revalid.ValidateStruct(usr)
	fmt.Println(err)
}
