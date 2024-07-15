package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /exercises", app.PostExercise)
	mux.HandleFunc("GET /exercises", app.GetExercises)
	mux.HandleFunc("PUT /exercises", app.PutExercise)
	mux.HandleFunc("DELETE /exercises", app.DeleteExercise)
	mux.HandleFunc("GET /exercises/{id}", app.GetExercise)
	mux.HandleFunc("GET /muscles", app.GetMuscles)

	chain := alice.New(app.logger)
	return chain.Then(mux)
}