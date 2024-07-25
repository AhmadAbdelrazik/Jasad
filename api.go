package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// handler signature
type apiFunc func(w http.ResponseWriter, r *http.Request) error

type apiErr struct {
	Error  string `json:"error"`
	status int 
}

// Error handling here
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			switch i := err.(type) {
			case ErrBadRequest:
				WriteJSON(w, i.Status, apiErr{Error: i.Message})
			case ErrNotFound:
				WriteJSON(w, i.Status, apiErr{Error: i.Message})
			default:
				WriteJSON(w, http.StatusBadRequest, apiErr{Error: i.Error()})
			}
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type APIServer struct {
	listenAddr string
	InfoLog    log.Logger
	ErrorLog   log.Logger
	DB         Storage
}

func NewAPIServer(listenAddr string, DB Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		InfoLog:    *log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate),
		ErrorLog:   *log.New(os.Stdout, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile),
		DB:         DB,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /exercises", makeHTTPHandleFunc(s.HandleGetExercises))
	mux.HandleFunc("GET /exercises/{id}", makeHTTPHandleFunc(s.HandleGetExerciseByID))

	mux.HandleFunc("POST /exercises", makeHTTPHandleFunc(s.HandleCreateExercise))

	http.ListenAndServe(s.listenAddr, mux)
}

func (s *APIServer) HandleCreateExercise(w http.ResponseWriter, r *http.Request) error {
	ExerciseRequest := CreateExerciseRequest{}
	if err := json.NewDecoder(r.Body).Decode(&ExerciseRequest); err != nil {
		return NewErrBadRequest()
	}

	// Validation (need to implement)
	if err := validate.Struct(ExerciseRequest); err != nil {
		return NewErrBadRequest()
	}

	if err := s.DB.CreateExercise(&ExerciseRequest); err != nil {

	}

	WriteJSON(w, http.StatusAccepted, ExerciseRequest)
	return nil
}

func (s *APIServer) HandleGetExercises(w http.ResponseWriter, r *http.Request) error {
	WriteJSON(w, http.StatusOK, &Exercise{})
	return nil
}

func (s *APIServer) HandleGetExerciseByID(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid id: %v", idStr)
	}

	WriteJSON(w, http.StatusOK, &Exercise{ExerciseID: id})
	return nil
}
