package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/delveper/mylib/app/models"
)

func main() {
	u, err := url.Parse("http://localhost:9999/books/download")
	if err != nil {
		log.Fatal(err)
	}

	f, err := models.NewDataFilter[models.Book](u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", f)
}
