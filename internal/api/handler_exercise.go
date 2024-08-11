package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
)

func (a *Application) HandleCreateExercise(w http.ResponseWriter, r *http.Request) {
	ExerciseRequest := storage.ExerciseCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&ExerciseRequest); err != nil {
		a.BadRequest(w)
	}
	r.Body.Close()

	if err := a.Validate.Struct(ExerciseRequest); err != nil {
		a.BadRequest(w)
	}

	if err := a.DB.CreateExercise(&ExerciseRequest); err != nil {
		if err == storage.ErrNoRecord {
			a.BadRequest(w)
		} else if strings.Contains(err.Error(), "Duplicate entry") {
			a.ClientError(w, http.StatusConflict)
		} else {
			a.ServerError(w, err)
		}
	}

	WriteJSON(w, http.StatusAccepted, APIResponse{Message: `exercise has been created`})

}

func (a *Application) HandleGetExercises(w http.ResponseWriter, r *http.Request) {

	exercises, err := a.DB.GetExercises()
	if err != nil {
		a.ServerError(w, err)
	}

	WriteJSON(w, http.StatusOK, exercises)

}

func (a *Application) HandleGetExercisesByMuscle(w http.ResponseWriter, r *http.Request) {
	var muscle storage.Muscle

	muscle.MuscleGroup = r.PathValue("muscleGroup")
	muscle.MuscleName = r.PathValue("muscleName")

	if err := a.DB.MuscleExists(&muscle); err != nil {
		a.BadRequest(w)
	}

	exercises, err := a.DB.GetExercisesByMuscle(muscle)
	if err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
	}

	WriteJSON(w, http.StatusOK, exercises)

}

func (a *Application) HandleGetExerciseByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		a.BadRequest(w)
	}

	exercise, err := a.DB.GetExerciseByID(id)
	if err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
	}

	WriteJSON(w, http.StatusOK, exercise)

}

func (a *Application) HandleGetExerciseByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	name = strings.ReplaceAll(name, "-", " ")

	exercise, err := a.DB.GetExerciseByName(name)
	if err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
	}

	WriteJSON(w, http.StatusOK, exercise)

}

func (a *Application) HandleUpdateExercise(w http.ResponseWriter, r *http.Request) {
	exercise := storage.ExerciseUpdateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		a.BadRequest(w)
	}
	r.Body.Close()

	if err := a.Validate.Struct(exercise); err != nil {
		a.BadRequest(w)
	}

	if err := a.DB.UpdateExercise(&exercise); err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
	}

	WriteJSON(w, http.StatusAccepted, APIResponse{Message: `exercise has been updated`})

}

func (a *Application) HandleDeleteExercise(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		a.BadRequest(w)
	}

	if err := a.DB.DeleteExercise(id); err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
	}

	WriteJSON(w, http.StatusAccepted, APIResponse{Message: `exercise has been deleted`})

}
