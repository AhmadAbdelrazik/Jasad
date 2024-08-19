package api

import (
	"net/http"

	"github.com/justinas/alice"
)

func (s *Application) Run() error {
	mux := http.NewServeMux()

	standard := alice.New(s.recoverPanic, secureHeaders, s.Logger)
	userInfoAuth := alice.New(s.Authenticate, s.AuthorizeUserInfo)

	// Users
	mux.HandleFunc("POST /signup", s.HandleSignup)
	mux.HandleFunc("POST /signin", s.HandleSignIn)
	mux.Handle("GET /user/{user}", userInfoAuth.Then(http.HandlerFunc(s.HandleUserInfo)))

	// Exercises
	mux.HandleFunc("POST /exercises", s.HandleCreateExercise)

	mux.HandleFunc("GET /exercises", s.HandleGetExercises)
	mux.HandleFunc("GET /exercises/muscle/{muscleGroup}/{muscleName}", s.HandleGetExercisesByMuscle)
	mux.HandleFunc("GET /exercises/id/{id}", s.HandleGetExerciseByID)
	mux.HandleFunc("GET /exercises/name/{name}", s.HandleGetExerciseByName)

	mux.HandleFunc("PUT /exercises", s.HandleUpdateExercise)

	mux.HandleFunc("DELETE /exercises/{id}", s.HandleDeleteExercise)

	// Workout Sessions

	// mux.HandleFunc("GET /user/{user}/workouts")
	// mux.HandleFunc("GET /user/{user}/workouts/{workout}")
	// mux.HandleFunc("GET /user/{user}/workouts/new")
	// mux.HandleFunc("POST /user/{user}/workouts")
	// mux.HandleFunc("PUT /user/{user}/workouts/{workout}")
	// mux.HandleFunc("DELETE /user/{user}/workouts/{workout}")

	err := http.ListenAndServe(s.Config.Port, standard.Then(mux))
	return err
}
