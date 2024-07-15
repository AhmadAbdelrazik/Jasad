package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *Application) BadRequestForm(w http.ResponseWriter, ErrorFields map[string]string) {
	errors := struct {
		Errors map[string]string `json:"errors"`
	}{
		Errors: ErrorFields,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusUnprocessableEntity)

	if err := json.NewEncoder(w).Encode(errors); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}