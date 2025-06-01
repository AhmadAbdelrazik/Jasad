package application

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/exercises", app.createExerciseHandler)
	mux.HandleFunc("GET /v1/exercises", app.searchExercisesHandler)
	mux.HandleFunc("GET /v1/exercises/{id}", app.getExericseHandler)
	mux.HandleFunc("PATCH /v1/exercises/{id}", app.updateExerciseHandler)
	mux.HandleFunc("DELETE /v1/exercises/{id}", app.deleteExerciseHandler)

	return mux
}
