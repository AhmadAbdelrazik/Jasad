package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
)

// HandleCreateExercise Creates a new Exercise
func (a *Application) HandleCreateExercise(w http.ResponseWriter, r *http.Request) {
	// Initialize ExerciseCreateRequest object to parse the request body
	ExerciseRequest := storage.ExerciseCreateRequest{}

	// Parse the request body using json, returns error at failure
	if err := json.NewDecoder(r.Body).Decode(&ExerciseRequest); err != nil {
		a.BadRequest(w)
		return
	}

	defer r.Body.Close()

	// Validate the response
	if err := a.Validate.Struct(ExerciseRequest); err != nil {
		a.BadRequest(w)
		return
	}

	// Database Call
	if err := a.DB.Exercise.CreateExercise(&ExerciseRequest); err != nil {
		switch err {
		case storage.ErrInvalidMuscle:
			a.BadRequest(w)
		case storage.ErrDuplicateEntry:
			a.ClientError(w, http.StatusConflict)
		default:
			a.ServerError(w, err)
		}
		return
	}

	// Response at success
	WriteJSON(w, http.StatusAccepted, APIResponse{Message: `exercise has been created`})
}

// HandleGetExercises Get all exercises in the database
func (a *Application) HandleGetExercises(w http.ResponseWriter, r *http.Request) {
	// Database Call
	exercises, err := a.DB.Exercise.GetExercises()
	if err != nil {
		a.ServerError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, exercises)
}

// HandleGetExercisesByMuscle Get all exercises that has a specific muscle
func (a *Application) HandleGetExercisesByMuscle(w http.ResponseWriter, r *http.Request) {
	var muscle storage.Muscle

	// Parse muscle name and muscle group from the path
	muscle.MuscleGroup = r.PathValue("muscleGroup")
	muscle.MuscleName = r.PathValue("muscleName")

	// Check if the muscle exists
	if err := a.DB.Exercise.MuscleExists(&muscle); err != nil {
		a.BadRequest(w)
		return
	}

	// Get all exercises that has the specific muscle
	exercises, err := a.DB.Exercise.GetExercisesByMuscle(muscle)
	if err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	WriteJSON(w, http.StatusOK, exercises)

}

// HandleGetExerciseByID Gets an exercise using it's id
func (a *Application) HandleGetExerciseByID(w http.ResponseWriter, r *http.Request) {
	// Get the id string from the path value
	idStr := r.PathValue("id")

	// Convert the string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		a.BadRequest(w)
		return
	}

	// Database call
	exercise, err := a.DB.Exercise.GetExerciseByID(id)
	if err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	WriteJSON(w, http.StatusOK, exercise)

}

// HandleGetExerciseByName Get the exercise using it's name
func (a *Application) HandleGetExerciseByName(w http.ResponseWriter, r *http.Request) {
	// Get Exercise name from the URL
	name := r.PathValue("name")

	// Replace '-' with ' '. 'front-shoulder' become 'front shoulder'
	name = strings.ReplaceAll(name, "-", " ")

	// Database Call
	exercise, err := a.DB.Exercise.GetExerciseByName(name)
	if err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	WriteJSON(w, http.StatusOK, exercise)
}

func (a *Application) HandleUpdateExercise(w http.ResponseWriter, r *http.Request) {
	// Get the id string from the path value
	idStr := r.PathValue("id")

	// Convert the string to int
	exerciseID, err := strconv.Atoi(idStr)
	if err != nil {
		a.BadRequest(w)
		return
	}

	// Initialize ExerciseUpdateRequest object to parse the request body
	exercise := storage.ExerciseUpdateRequest{}

	// Parse the request body using json, returns error at failure
	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		a.BadRequest(w)
		return
	}

	defer r.Body.Close()

	// Validate the response
	if err := a.Validate.Struct(exercise); err != nil {
		a.BadRequest(w)
	}

	// Database Call
	if err := a.DB.Exercise.UpdateExercise(exerciseID, &exercise); err != nil {
		switch err {
		case storage.ErrInvalidMuscle:
			a.BadRequest(w)
		case storage.ErrNoRecord:
			a.NotFound(w)
		default:
			a.ServerError(w, err)
		}
		return
	}

	WriteJSON(w, http.StatusAccepted, APIResponse{Message: `exercise has been updated`})

}

func (a *Application) HandleDeleteExercise(w http.ResponseWriter, r *http.Request) {
	// Get the id string from the path value
	idStr := r.PathValue("id")

	// Convert the string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		a.BadRequest(w)
		return
	}

	// Database call
	if err := a.DB.Exercise.DeleteExercise(id); err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	WriteJSON(w, http.StatusAccepted, APIResponse{Message: `exercise has been deleted`})

}
