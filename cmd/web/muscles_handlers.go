package main

import (
	"encoding/json"
	"net/http"

	"github.com/AhmadAbdelrazik/jasad/internal/model"
)

func (app *Application) GetMuscles(w http.ResponseWriter, r *http.Request) {
	muscles, err := app.Muscle.GetAll()
	if err != nil {
		if err == model.ErrNoRecord{
			app.ClientError(w, http.StatusBadRequest)
		} else {
			app.ServerError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(muscles); err != nil {
		app.ServerError(w, err)
		return
	}
}