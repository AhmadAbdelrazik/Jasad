package api

import (
	"net/http"

	"github.com/justinas/alice"
)

func (a *Application) Run() error {
	mux := http.NewServeMux()

	standard := alice.New(a.recoverPanic, secureHeaders, a.Logger)
	userAuth := alice.New(a.Authenticate, a.AuthorizeUserInfo)

	// Users
	mux.HandleFunc("POST /signup", a.HandleSignup)
	mux.HandleFunc("POST /signin", a.HandleSignIn)
	mux.Handle("GET /user/{user}", userAuth.ThenFunc(a.HandleUserInfo))

	// Exercises
	mux.HandleFunc("POST /exercises", a.HandleCreateExercise)

	mux.HandleFunc("GET /exercises", a.HandleGetExercises)
	mux.HandleFunc("GET /exercises/muscle/{muscleGroup}/{muscleName}", a.HandleGetExercisesByMuscle)
	mux.HandleFunc("GET /exercises/id/{id}", a.HandleGetExerciseByID)
	mux.HandleFunc("GET /exercises/name/{name}", a.HandleGetExerciseByName)

	mux.HandleFunc("PUT /exercises", a.HandleUpdateExercise)

	mux.HandleFunc("DELETE /exercises/{id}", a.HandleDeleteExercise)

	// Workout Sessions
	mux.Handle("GET /user/{user}/workouts", userAuth.ThenFunc(a.HandleGetWorkouts))
	mux.Handle("GET /user/{user}/workouts/new", userAuth.ThenFunc(a.HandleCreateWorkoutForm))
	mux.Handle("POST /user/{user}/workouts/new", userAuth.ThenFunc(a.HandleCreateWorkout))

	mux.Handle("GET /user/{user}/workouts/{workout}", userAuth.ThenFunc(a.HandleGetWorkout))
	mux.Handle("PUT /user/{user}/workouts/{workout}", userAuth.ThenFunc(a.HandleUpdateWorkout))
	mux.Handle("DELETE /user/{user}/workouts/{workout}", userAuth.ThenFunc(a.HandleDeleteWorkout))

	err := http.ListenAndServe(a.Config.Port, standard.Then(mux))
	return err
}
