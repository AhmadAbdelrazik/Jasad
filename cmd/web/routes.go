package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /exercises", app.exerciseCreate)
	mux.HandleFunc("GET /exercises", app.exerciseRead)
	mux.HandleFunc("PUT /exercises", app.exerciseUpdate)
	mux.HandleFunc("DELETE /exercises", app.exerciseDelete)
	mux.HandleFunc("GET /exercises/{id}", app.exerciseID)
	mux.HandleFunc("GET /muscles", app.musclesReadAll)

	chain := alice.New(app.logger)
	return chain.Then(mux)
}