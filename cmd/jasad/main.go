package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile)

	cache := cache.NewRedis()

	app := &api.Application{
		Config:   config,
		DB:       db,
		Cache:    cache,
		Validate: validate,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	server := &http.Server{
		Addr:         config.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      app.Routes(),
		ErrorLog:     errorLog,
	}

	app.InfoLog.Println(fmt.Sprint("Started listening at port ", config.Port))
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
