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

// handler signature
type apiFunc func(w http.ResponseWriter, r *http.Request) error

type apiErr struct {
	Message string `json:"error"`
	Status  int    `json:"status"`
}

type apiResponse struct {
	Message string `json:"message"`
}

func (r *apiErr) Error() string {
	return r.Message
}

// Error handling here
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		switch i := err.(type) {
		case *apiErr:
			WriteJSON(w, i.Status, i)
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type APIServer struct {
	config     *config.Configuration
	listenAddr string
	InfoLog    log.Logger
	ErrorLog   log.Logger
	DB         storage.Storage
	Validate   *validator.Validate
}

func NewAPIServer(listenAddr string, DB storage.Storage, validate *validator.Validate) *APIServer {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	return &APIServer{
		config:     config,
		listenAddr: listenAddr,
		InfoLog:    *log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate),
		ErrorLog:   *log.New(os.Stdout, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile),
		DB:         DB,
		Validate:   validate,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()

	log := alice.New(s.logger)

	mux.HandleFunc("POST /signup", makeHTTPHandleFunc(s.HandleSignup))
	mux.HandleFunc("POST /signin", makeHTTPHandleFunc(s.HandleSignIn))

	mux.HandleFunc("POST /exercises", makeHTTPHandleFunc(s.HandleCreateExercise))

	mux.HandleFunc("GET /exercises", makeHTTPHandleFunc(s.HandleGetExercises))
	mux.HandleFunc("GET /exercises/muscle/{muscleGroup}/{muscleName}", makeHTTPHandleFunc(s.HandleGetExercisesByMuscle))
	mux.HandleFunc("GET /exercises/id/{id}", makeHTTPHandleFunc(s.HandleGetExerciseByID))
	mux.HandleFunc("GET /exercises/name/{name}", makeHTTPHandleFunc(s.HandleGetExerciseByName))

	mux.HandleFunc("PUT /exercises", makeHTTPHandleFunc(s.HandleUpdateExercise))

	mux.HandleFunc("DELETE /exercises/{id}", makeHTTPHandleFunc(s.HandleDeleteExercise))

	http.ListenAndServe(s.listenAddr, log.Then(mux))
}
