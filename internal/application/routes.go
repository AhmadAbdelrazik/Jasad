package application

import (
	"net/http"

	"github.com/ahmadabdelrazik/jasad/internal/model"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/exercises", app.IsAuthorized(app.createExerciseHandler))
	mux.HandleFunc("GET /v1/exercises", app.searchExercisesHandler)
	mux.HandleFunc("GET /v1/exercises/{id}", app.getExericseHandler)
	mux.HandleFunc("PATCH /v1/exercises/{id}", app.IsAuthorized(app.updateExerciseHandler))
	mux.HandleFunc("DELETE /v1/exercises/{id}", app.IsAuthorized(app.deleteExerciseHandler))

	mux.HandleFunc("GET /google_login", app.googleLoginHandler)
	mux.HandleFunc("GET /google_callback", app.googleCallbackHandler)

	mux.HandleFunc("GET /v1/users", app.IsAuthorized(app.GetAllUsers))
	mux.HandleFunc("GET /v1/users/{id}", app.IsAuthorized(app.getUserByIDHandler, model.RoleUser))

	return app.recoverPanic(app.rateLimit(mux))
}
