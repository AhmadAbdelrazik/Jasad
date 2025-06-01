package application

import (
	"errors"
	"net/http"

	"github.com/ahmadabdelrazik/jasad/internal/model"
	"github.com/ahmadabdelrazik/jasad/pkg/validator"
)

func (app *Application) createExerciseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string `json:"name"`
		Muscle         string `json:"muscle"`
		Instructions   string `json:"instructions"`
		AdditionalInfo string `json:"additional_info"`
		ImageURL       string `json:"image_url"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	muscle, err := model.GetMuscle(input.Muscle)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	exercise := &model.Exercise{
		Name:           input.Name,
		Muscle:         muscle,
		Instructions:   input.Instructions,
		AdditionalInfo: input.AdditionalInfo,
		ImageURL:       input.ImageURL,
	}

	v := validator.New()
	exercise.Validate(v)
	if !v.Valid() {
		FailedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Exercises.Create(exercise); err != nil {
		switch {
		case errors.Is(err, model.ErrAlreadyExists):
			ConflictResponse(w, r)
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"exercise": exercise}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func (app *Application) getExericseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	exercise, err := app.models.Exercises.Get(int(id))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			NotFoundResponse(w, r)
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"exercise": exercise}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func (app *Application) searchExercisesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   string `json:"name"`
		Muscle string `json:"muscle"`
		model.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Muscle = app.readString(qs, "muscle", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafeList = []string{"muscle", "id", "name", "-muscle", "-id", "-name"}

	model.ValidateFilters(v, input.Filters)

	if !v.Valid() {
		FailedValidationResponse(w, r, v.Errors)
		return
	}

	exercises, metadata, err := app.models.Exercises.Search(input.Name, input.Muscle, input.Filters)
	if err != nil {
		ServerErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"exercises": exercises, "metadata": metadata}, nil)
}

func (app *Application) updateExerciseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	var input struct {
		Name           *string `json:"name"`
		Muscle         *string `json:"muscle"`
		Instructions   *string `json:"instructions"`
		AdditionalInfo *string `json:"additional_info"`
		ImageURL       *string `json:"image_url"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	exercise, err := app.models.Exercises.Get(int(id))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			NotFoundResponse(w, r)
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	if input.Name != nil {
		exercise.Name = *input.Name
	}
	if input.Muscle != nil {
		muscle, err := model.GetMuscle(*input.Muscle)
		if err != nil {
			BadRequestResponse(w, r, err)
			return
		}
		exercise.Muscle = muscle
	}
	if input.Instructions != nil {
		exercise.Instructions = *input.Instructions
	}
	if input.AdditionalInfo != nil {
		exercise.AdditionalInfo = *input.AdditionalInfo
	}
	if input.ImageURL != nil {
		exercise.ImageURL = *input.ImageURL
	}

	v := validator.New()
	exercise.Validate(v)

	if !v.Valid() {
		FailedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Exercises.Update(exercise); err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			EditConflictResponse(w, r)
		case errors.Is(err, model.ErrAlreadyExists):
			ConflictResponse(w, r)
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(
		w,
		http.StatusOK,
		envelope{"message": "updated successfully", "exercise": exercise},
		nil,
	)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func (app *Application) deleteExerciseHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	if err := app.models.Exercises.Delete(int(id)); err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			NotFoundResponse(w, r)
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "exercise deleted successfully"}, nil)
}
