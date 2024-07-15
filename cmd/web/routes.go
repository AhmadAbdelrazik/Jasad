package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	// Exercises Handler
	mux.HandleFunc("POST /exercises", app.PostExercise)
	mux.HandleFunc("GET /exercises", app.GetExercises)
	mux.HandleFunc("PUT /exercises/{id}", app.PutExercise)
	mux.HandleFunc("DELETE /exercises/{id}", app.DeleteExercise)
	mux.HandleFunc("GET /exercises/{id}", app.GetExercise)
	
	// Muscles Handler
	mux.HandleFunc("GET /muscles", app.GetMuscles)

	chain := alice.New(app.logger)
	return chain.Then(mux)
}