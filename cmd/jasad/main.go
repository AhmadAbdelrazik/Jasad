package main

import (
	"fmt"
	"log"

	"github.com/AhmadAbdelrazik/jasad/internal/api"
	"github.com/AhmadAbdelrazik/jasad/internal/cache"
	"github.com/AhmadAbdelrazik/jasad/internal/config"
	"github.com/AhmadAbdelrazik/jasad/internal/storage"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func main() {
	validate = validator.New()

	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := storage.NewMySQLDatabase(config.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	cache := cache.NewRedis()

	server := api.NewApplication(config, db, cache, validate)
	server.InfoLog.Println(fmt.Sprint("Started listening at port ", config.Port))
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
