package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
)

func (a *Application) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workoutRequest storage.WorkoutCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&workoutRequest); err != nil {
		a.BadRequest(w)
		return
	}

	if err := a.Validate.Struct(workoutRequest); err != nil {
		a.BadRequest(w)
		return
	}

	userID := r.Context().Value("userID").(int)
	sessionID, err := a.DB.Workout.CreateWorkout(workoutRequest, userID)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d/workouts/%d", userID, sessionID), http.StatusPermanentRedirect)
}

func (a *Application) HandleCreateWorkoutForm(w http.ResponseWriter, r *http.Request) {
	// Send all Information about exercises.
	exercises, err := a.DB.Exercise.GetExercises()
	if err != nil {
		a.ServerError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, exercises)
}

func (a *Application) HandleGetWorkout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.Atoi(r.PathValue("workout"))
	if err != nil {
		a.BadRequest(w)
		return
	}

	userID := r.Context().Value("userID").(int)

	session, err := a.DB.Workout.GetWorkout(sessionID, userID)
	if err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	var exercises []storage.Exercise
	for _, workout := range session.Workouts {
		e, err := a.DB.Exercise.GetExerciseByID(workout.ExerciseID)
		if err != nil {
			if err == storage.ErrNoRecord {
				a.BadRequest(w)
			} else {
				a.ServerError(w, err)
			}
			return
		}
		exercises = append(exercises, *e)

	}

	response := struct {
		Session   storage.Session    `json:"session"`
		Exercises []storage.Exercise `json:"exercises"`
	}{
		Session:   *session,
		Exercises: exercises,
	}

	WriteJSON(w, http.StatusAccepted, response)
}

func (a *Application) HandleGetWorkouts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	username := r.PathValue("user")

	sessions, err := a.DB.Workout.GetWorkouts(userID)
	if err != nil && err != storage.ErrNoRecord {
		a.ServerError(w, err)
		return
	}

	res := &storage.WorkoutGetResponse{
		Username: username,
		Sessions: sessions,
	}

	WriteJSON(w, http.StatusOK, res)
}

func (a *Application) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {

}

func (a *Application) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {

}
