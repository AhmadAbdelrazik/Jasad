package main

import (
	"log"

	"github.com/AhmadAbdelrazik/jasad/internal/api"
	"github.com/AhmadAbdelrazik/jasad/internal/storage"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func main() {
	validate = validator.New()

	db, err := storage.NewMySQLDatabase()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":3000", db, validate)
	server.InfoLog.Println("Server is listening on port 3000")
	server.Run()
}
