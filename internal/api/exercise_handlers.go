package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
)

func (s *APIServer) HandleCreateExercise(w http.ResponseWriter, r *http.Request) error {
	ExerciseRequest := storage.ExerciseCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&ExerciseRequest); err != nil {
		return s.BadRequest()
	}
	r.Body.Close()

	if err := s.Validate.Struct(ExerciseRequest); err != nil {
		return s.BadRequest()
	}

	if err := s.DB.CreateExercise(&ExerciseRequest); err != nil {
		if err == storage.ErrNoRecord {
			return s.BadRequest()
		} else if strings.Contains(err.Error(), "Duplicate entry") {
			return s.ClientError(http.StatusConflict)
		} else {
			return s.ServerError(err)
		}
	}

	WriteJSON(w, http.StatusAccepted, apiResponse{Message: `exercise has been created`})
	return nil
}

func (s *APIServer) HandleGetExercises(w http.ResponseWriter, r *http.Request) error {

	exercises, err := s.DB.GetExercises()
	if err != nil {
		return s.ServerError(err)
	}

	WriteJSON(w, http.StatusOK, exercises)
	return nil
}

func (s *APIServer) HandleGetExercisesByMuscle(w http.ResponseWriter, r *http.Request) error {
	var muscle storage.Muscle

	muscle.MuscleGroup = r.PathValue("muscleGroup")
	muscle.MuscleName = r.PathValue("muscleName")

	if err := s.DB.MuscleExists(&muscle); err != nil {
		return s.BadRequest()
	}

	exercises, err := s.DB.GetExercisesByMuscle(muscle)
	if err != nil {
		if err == storage.ErrNoRecord {
			return s.NotFound()
		} else {
			return s.ServerError(err)
		}
	}

	WriteJSON(w, http.StatusOK, exercises)

	return nil
}

func (s *APIServer) HandleGetExerciseByID(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return s.BadRequest()
	}

	exercise, err := s.DB.GetExerciseByID(id)
	if err != nil {
		if err == storage.ErrNoRecord {
			return s.NotFound()
		} else {
			return s.ServerError(err)
		}
	}

	WriteJSON(w, http.StatusOK, exercise)
	return nil
}

func (s *APIServer) HandleGetExerciseByName(w http.ResponseWriter, r *http.Request) error {
	name := r.PathValue("name")

	name = strings.ReplaceAll(name, "-", " ")

	exercise, err := s.DB.GetExerciseByName(name)
	if err != nil {
		if err == storage.ErrNoRecord {
			return s.NotFound()
		} else {
			return s.ServerError(err)
		}
	}

	WriteJSON(w, http.StatusOK, exercise)
	return nil
}

func (s *APIServer) HandleUpdateExercise(w http.ResponseWriter, r *http.Request) error {
	exercise := storage.ExerciseUpdateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		return s.BadRequest()
	}
	r.Body.Close()

	if err := s.Validate.Struct(exercise); err != nil {
		return s.BadRequest()
	}

	if err := s.DB.UpdateExercise(&exercise); err != nil {
		if err == storage.ErrNoRecord {
			return s.NotFound()
		} else {
			return s.ServerError(err)
		}
	}

	WriteJSON(w, http.StatusAccepted, apiResponse{Message: `exercise has been updated`})
	return nil
}

func (s *APIServer) HandleDeleteExercise(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return s.BadRequest()
	}

	if err := s.DB.DeleteExercise(id); err != nil {
		if err == storage.ErrNoRecord {
			return s.NotFound()
		} else {
			return s.ServerError(err)
		}
	}

	WriteJSON(w, http.StatusAccepted, apiResponse{Message: `exercise has been deleted`})
	return nil
}
