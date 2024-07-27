package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/justinas/alice"
)

// handler signature
type apiFunc func(w http.ResponseWriter, r *http.Request) error

type apiErr struct {
	Message string `json:"error"`
	Status  int    `json:"status"`
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

	mux.HandleFunc("POST /exercises", makeHTTPHandleFunc(s.HandleCreateExercise))
	mux.HandleFunc("GET /exercises", makeHTTPHandleFunc(s.HandleGetExercises))
	mux.HandleFunc("GET /exercises/{id}", makeHTTPHandleFunc(s.HandleGetExerciseByID))
	mux.HandleFunc("PUT /exercises", makeHTTPHandleFunc(s.HandleUpdateExercise))

	log := alice.New(s.logger)

	http.ListenAndServe(s.listenAddr, log.Then(mux))
}

func (s *APIServer) HandleCreateExercise(w http.ResponseWriter, r *http.Request) error {
	ExerciseRequest := CreateExerciseRequest{}
	if err := json.NewDecoder(r.Body).Decode(&ExerciseRequest); err != nil {
		return s.BadRequest()
	}
	r.Body.Close()

	if err := validate.Struct(ExerciseRequest); err != nil {
		return s.BadRequest()
	}

	if err := s.DB.CreateExercise(&ExerciseRequest); err != nil {
		if err == ErrNoRecord {
			return s.BadRequest()
		} else if strings.Contains(err.Error(), "Duplicate entry") {
			return s.ClientError(http.StatusConflict)
		} else {
			s.ServerError(err)
		}
	}

	WriteJSON(w, http.StatusAccepted, ExerciseRequest)
	return nil
}

func (s *APIServer) HandleGetExercises(w http.ResponseWriter, r *http.Request) error {

	val, err := s.DB.GetExercises()
	if err != nil {
		return s.ServerError(err)
	}

	WriteJSON(w, http.StatusOK, val)
	return nil
}

func (s *APIServer) HandleGetExerciseByID(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.BadRequest()
	}

	WriteJSON(w, http.StatusOK, &Exercise{ExerciseID: id})
	return nil
}

func (s *APIServer) HandleUpdateExercise(w http.ResponseWriter, r *http.Request) error {
	exercise := UpdateExerciseRequest{}

	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		return s.BadRequest()
	}
	r.Body.Close()

	if err := validate.Struct(exercise); err != nil {
		return s.BadRequest()
	}

	if err := s.DB.UpdateExercise(&exercise); err != nil {
		if err == ErrNoRecord {
			return s.NotFound()
		} else {
			return s.ServerError(err)
		}
	}

	WriteJSON(w, http.StatusAccepted, exercise)
	return nil
}
