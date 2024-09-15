package api

import (
	"net/http"

	"github.com/justinas/alice"
)

// Routes Contains the main application multiplexer.
// Here We specify all the end points patterns
func (a *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	standard := alice.New(a.recoverPanic, secureHeaders, a.Logger, a.RateLimiter)
	userAuth := alice.New(a.Authenticate, a.AuthorizeUserInfo)
	adminAuth := alice.New(a.Authenticate, a.AuthorizeAdmin)

	// Users
	mux.HandleFunc("POST /api/v1/users/signup", a.HandleSignup)
	mux.HandleFunc("POST /api/v1/users/signin", a.HandleSignIn)
	mux.Handle("GET /api/v1/users/{user}/info", adminAuth.ThenFunc(a.HandleUserInfo))

	// Exercises

	mux.HandleFunc("GET /api/v1/exercises", a.HandleGetExercises)
	mux.HandleFunc("GET /api/v1/exercises/muscle/{muscleGroup}/{muscleName}", a.HandleGetExercisesByMuscle)
	mux.HandleFunc("GET /api/v1/exercises/id/{id}", a.HandleGetExerciseByID)
	mux.HandleFunc("GET /api/v1/exercises/name/{name}", a.HandleGetExerciseByName)

	mux.Handle("POST /api/v1/exercises", adminAuth.ThenFunc(a.HandleCreateExercise))
	mux.Handle("PUT /api/v1/exercises/{id}", adminAuth.ThenFunc(a.HandleUpdateExercise))
	mux.Handle("DELETE /api/v1/exercises/{id}", adminAuth.ThenFunc(a.HandleDeleteExercise))

	// Workout Sessions
	mux.Handle("GET /api/v1/users/{user}/workouts", userAuth.ThenFunc(a.HandleGetWorkouts))
	mux.Handle("GET /api/v1/users/{user}/workouts/new", userAuth.ThenFunc(a.HandleCreateWorkoutForm))
	mux.Handle("POST /api/v1/users/{user}/workouts/new", userAuth.ThenFunc(a.HandleCreateWorkout))

	mux.Handle("GET /api/v1/users/{user}/workouts/{workout}", userAuth.ThenFunc(a.HandleGetWorkout))
	mux.Handle("PUT /api/v1/users/{user}/workouts/{workout}", userAuth.ThenFunc(a.HandleUpdateWorkout))
	mux.Handle("DELETE /api/v1/users/{user}/workouts/{workout}", userAuth.ThenFunc(a.HandleDeleteWorkout))

	return standard.Then(mux)
}
