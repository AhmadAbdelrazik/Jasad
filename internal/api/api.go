package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/AhmadAbdelrazik/jasad/internal/config"
	"github.com/AhmadAbdelrazik/jasad/internal/storage"
	"github.com/go-playground/validator/v10"
	"github.com/justinas/alice"
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
	DB       storage.Storage
	Validate *validator.Validate
}

func NewApplication(
	config *config.Configuration,
	DB storage.Storage,
	validate *validator.Validate) *Application {
	return &Application{
		Config:   config,
		InfoLog:  *log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate),
		ErrorLog: *log.New(os.Stdout, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile),
		DB:       DB,
		Validate: validate,
	}
}

func (s *Application) Run() {
	mux := http.NewServeMux()

	log := alice.New(s.Logger)

	mux.HandleFunc("POST /signup", s.HandleSignup)
	mux.HandleFunc("POST /signin", s.HandleSignIn)
	mux.HandleFunc("POST /exercises", s.HandleCreateExercise)
	mux.HandleFunc("GET /exercises", s.HandleGetExercises)
	mux.HandleFunc("GET /exercises/muscle/{muscleGroup}/{muscleName}", s.HandleGetExercisesByMuscle)
	mux.HandleFunc("GET /exercises/id/{id}", s.HandleGetExerciseByID)
	mux.HandleFunc("GET /exercises/name/{name}", s.HandleGetExerciseByName)
	mux.HandleFunc("PUT /exercises", s.HandleUpdateExercise)
	mux.HandleFunc("DELETE /exercises/{id}", s.HandleDeleteExercise)

	http.ListenAndServe(s.Config.Port, log.Then(mux))
}
