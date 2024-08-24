package api

import (
	"net/http"

	"github.com/justinas/alice"
)

func (a *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	standard := alice.New(a.recoverPanic, secureHeaders, a.Logger)
	userAuth := alice.New(a.Authenticate, a.AuthorizeUserInfo)
	adminAuth := alice.New(a.Authenticate, a.AuthorizeAdmin)

	// Users
	mux.HandleFunc("POST /users/signup", a.HandleSignup)
	mux.HandleFunc("POST /users/signin", a.HandleSignIn)
	mux.Handle("GET /users/{user}/info", adminAuth.ThenFunc(a.HandleUserInfo))

	// Exercises

	mux.HandleFunc("GET /exercises", a.HandleGetExercises)
	mux.HandleFunc("GET /exercises/muscle/{muscleGroup}/{muscleName}", a.HandleGetExercisesByMuscle)
	mux.HandleFunc("GET /exercises/id/{id}", a.HandleGetExerciseByID)
	mux.HandleFunc("GET /exercises/name/{name}", a.HandleGetExerciseByName)

	mux.Handle("POST /exercises", adminAuth.ThenFunc(a.HandleCreateExercise))
	mux.Handle("PUT /exercises", adminAuth.ThenFunc(a.HandleUpdateExercise))
	mux.Handle("DELETE /exercises/{id}", adminAuth.ThenFunc(a.HandleDeleteExercise))

	// Workout Sessions
	mux.Handle("GET /users/{user}/workouts", userAuth.ThenFunc(a.HandleGetWorkouts))
	mux.Handle("GET /users/{user}/workouts/new", userAuth.ThenFunc(a.HandleCreateWorkoutForm))
	mux.Handle("POST /users/{user}/workouts/new", userAuth.ThenFunc(a.HandleCreateWorkout))

	mux.Handle("GET /users/{user}/workouts/{workout}", userAuth.ThenFunc(a.HandleGetWorkout))
	mux.Handle("PUT /users/{user}/workouts/{workout}", userAuth.ThenFunc(a.HandleUpdateWorkout))
	mux.Handle("DELETE /users/{user}/workouts/{workout}", userAuth.ThenFunc(a.HandleDeleteWorkout))

	return standard.Then(mux)
}
