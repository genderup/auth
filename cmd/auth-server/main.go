package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/inappcloud/auth"
	"github.com/zenazn/goji"
)

func main() {
	url := os.Getenv("DATABASE_URL")

	if len(url) == 0 {
		log.Fatal(errors.New("You must set DATABASE_URL environment variable."))
	}

	if len(os.Getenv("PRIVATE_KEY")) == 0 {
		log.Fatal(errors.New("You must set PRIVATE_KEY environment variable."))
	}

	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal(err)
	}

	goji.Handle("*", auth.Mux(db))
	goji.Serve()
}
