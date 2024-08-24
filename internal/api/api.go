package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AhmadAbdelrazik/jasad/internal/cache"
	"github.com/AhmadAbdelrazik/jasad/internal/config"
	"github.com/AhmadAbdelrazik/jasad/internal/storage"
	"github.com/go-playground/validator/v10"
)

// Writes a json response. Used at the end of a handler
// to send the response.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Simple API response, Used to send simple messages
type APIResponse struct {
	Message string `json:"message"`
}

type Application struct {
	Config   *config.Configuration
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DB       *storage.Storage
	Validate *validator.Validate
	Cache    cache.Cache
}
