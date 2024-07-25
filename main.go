package main

import (
	"log"
	"github.com/go-playground/validator/v10"
)


var validate *validator.Validate

func main() {
	validate = validator.New()
	
	db, err := NewMySQLServer()
	if err != nil {
		log.Fatal(err)
	}
	

	server := NewAPIServer(":3000", db)
	server.InfoLog.Println("Server is listening on port 3000")
	server.Run()
}