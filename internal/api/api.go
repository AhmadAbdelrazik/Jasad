package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/AhmadAbdelrazik/jasad/internal/config"
	"github.com/AhmadAbdelrazik/jasad/internal/storage"
	"github.com/go-playground/validator/v10"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type APIResponse struct {
	Message string `json:"message"`
}

type Application struct {
	Config   *config.Configuration
	InfoLog  log.Logger
	ErrorLog log.Logger
	DB       *storage.Storage
	Validate *validator.Validate
}

func NewApplication(
	config *config.Configuration,
	DB *storage.Storage,
	validate *validator.Validate) *Application {
	return &Application{
		Config:   config,
		DB:       DB,
		Validate: validate,
		InfoLog:  *log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate),
		ErrorLog: *log.New(os.Stdout, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile),
	}
}
