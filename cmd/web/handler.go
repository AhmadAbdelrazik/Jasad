package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AhmadAbdelrazik/jasad/internal/model"
	"github.com/AhmadAbdelrazik/jasad/internal/validator"
)

type ExerciseCreateForm struct {
	ExerciseName        string   `json:"exercise_name"`
	ExerciseDescription string   `json:"exercise_explanation"`
	ReferenceVideo      string   `json:"reference_video"`
	Muscles             []string `json:"muscles"`
	validator.Validator
}

func (app *Application) musclesReadAll(w http.ResponseWriter, r *http.Request) {
	muscles, err := app.Jasad.GetAllMuscles()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.ClientError(w, http.StatusBadRequest)
		} else {
			app.ServerError(w, err)
		}
		return
	}

	buffer, err := json.Marshal(muscles)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(buffer))
}

func (app *Application) exerciseCreate(w http.ResponseWriter, r *http.Request) {
	// parse the body
	var form ExerciseCreateForm
	err := json.NewDecoder(r.Body).Decode(&form)

	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// pass the form, get boolean
	CheckExerciseCreateForm(&form)

	if muscleName, err := app.Jasad.CheckMusclesExist(form.Muscles); err != nil {
		if err == model.ErrNoRecord {
			form.AddFieldError(muscleName, "muscle doesn't exist")
		} else {
			app.ServerError(w, err)
			return
		}
	}

	if !form.Valid() {
		app.BadRequestForm(w, form.FieldErrors)
		return
	}

	// add the body to the database
	id, err := app.Jasad.AddExercise(form.ExerciseName, form.ExerciseDescription, form.ReferenceVideo, form.Muscles)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	response := struct {
		Message string `json:"message"`
		Id      int    `json:"id"`
	} {
		Message: "success",
		Id: id,
	}

	// return success or failure
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		app.ServerError(w, err)
		return
	}
}

func (app *Application) exerciseRead(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) exerciseUpdate(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) exerciseDelete(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) exerciseID(w http.ResponseWriter, r *http.Request) {

}
