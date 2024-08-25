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

// Validator is used to validate structs.
// It is used to validate request bodies.
var validate *validator.Validate

func main() {
	validate = validator.New()

	// Load the configs from configs/development.yml
	config, err := config.NewConfig("development")
	if err != nil {
		log.Fatal(err)
	}

	// Load new mysql database
	db, err := storage.NewMySQLDatabase(config.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile)

	// Load Redis Cache
	cache := cache.NewRedis()

	// Initialize Application instance
	app := &api.Application{
		Config:   config,
		DB:       db,
		Cache:    cache,
		Validate: validate,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}

	// Initialize a new server, passing application routes to handler
	server := &http.Server{
		Addr:         config.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      app.Routes(),
		ErrorLog:     errorLog,
	}

	// Start the server
	app.InfoLog.Println(fmt.Sprint("Started listening at port ", config.Port))
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
