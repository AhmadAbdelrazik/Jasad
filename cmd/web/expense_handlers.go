package main

import (
	"encoding/json"
	"net/http"

	"github.com/AhmadAbdelrazik/jasad/internal/model"
	"github.com/AhmadAbdelrazik/jasad/internal/validator"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type ExerciseCreateForm struct {
	ExerciseName        string   `json:"exercise_name"`
	ExerciseDescription string   `json:"exercise_explanation"`
	ReferenceVideo      string   `json:"reference_video"`
	Muscles             []string `json:"muscles"`
	validator.Validator
}

func (app *Application) PostExercise(w http.ResponseWriter, r *http.Request) {
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

	if muscleName, err := app.Muscle.AllExist(form.Muscles); err != nil {
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
	id, err := app.Exercise.Create(form.ExerciseName, form.ExerciseDescription, form.ReferenceVideo, form.Muscles)
	if err != nil {
		ErrMysql, ok := err.(*mysql.MySQLError)
		if !ok {
			app.ServerError(w, err)
			return
		}
		if ErrMysql.Number == 1062 {
			app.ClientError(w, http.StatusConflict)
		}
		app.ServerError(w, err)
		return
	}

	response := struct {
		Message string    `json:"message"`
		Id      uuid.UUID `json:"id"`
	}{
		Message: "success",
		Id:      id,
	}

	// return success or failure
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		app.ServerError(w, err)
		return
	}
}

func (app *Application) GetExercise(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	Exercise, err := app.Exercise.GetByID(id)
	if err != nil {
		if err == model.ErrNoRecord {
			app.NotFound(w)
		} else {
			app.ServerError(w, err)
		}
		return
	}

	response := struct {
		ExerciseName        string   `json:"exercise_name"`
		ExerciseDescription string   `json:"exercise_explanation"`
		ReferenceVideo      string   `json:"reference_video"`
		Muscles             []string `json:"muscles"`
	}{
		ExerciseName: Exercise.ExerciseName,
		ExerciseDescription: Exercise.ExerciseDescription,
		ReferenceVideo: Exercise.ReferenceVideo,
		Muscles: Exercise.Muscles,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.ServerError(w, err)
	}
}

func (app *Application) PutExercise(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) DeleteExercise(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) GetExercises(w http.ResponseWriter, r *http.Request) {
	muscles, err := app.Exercise.GetAll()
	if err != nil {
		if err == model.ErrNoRecord {
			message := struct {
				Message string `json:"message"`
			}{Message: "No Exercises found"}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(message); err != nil {
				app.ServerError(w, err)
			}
			return
		} else {
			app.ServerError(w, err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(muscles); err != nil {
		app.ServerError(w, err)
	}
}
