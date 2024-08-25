package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
)

// HandleCreateWorkout Creates workout and link it to the user
func (a *Application) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	// Initialize WorkoutRequest object to parse the request body
	var workoutRequest storage.WorkoutRequest

	// Parse the request body using json, returns error at failure
	if err := json.NewDecoder(r.Body).Decode(&workoutRequest); err != nil {
		a.BadRequest(w)
		return
	}

	defer r.Body.Close()

	// Validate the response
	if err := a.Validate.Struct(workoutRequest); err != nil {
		a.BadRequest(w)
		return
	}

	// get the userID. userID is parsed in the auth middleware
	userID := r.Context().Value("userID").(int)

	// Database Call
	sessionID, err := a.DB.Workout.CreateWorkout(workoutRequest, userID)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	// Redirect to the newly created workout
	http.Redirect(w, r, fmt.Sprintf("/users/%d/workouts/%d", userID, sessionID), http.StatusPermanentRedirect)
}

// HandleCreateWorkoutForm Passes all the available exercises
// for the user to construct the workout
func (a *Application) HandleCreateWorkoutForm(w http.ResponseWriter, r *http.Request) {
	// Send all Information about exercises.
	exercises, err := a.DB.Exercise.GetExercises()
	if err != nil {
		a.ServerError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, exercises)
}

// HandleGetWorkout Gets a specific workout
func (a *Application) HandleGetWorkout(w http.ResponseWriter, r *http.Request) {
	// parse the sessionID value from the header and convert to number
	sessionID, err := strconv.Atoi(r.PathValue("workout"))
	if err != nil { // return error if the string conversion failed
		a.BadRequest(w)
		return
	}

	// get the userID. userID is parsed in the auth middleware
	userID := r.Context().Value("userID").(int)

	// Database Call
	session, err := a.DB.Workout.GetWorkout(sessionID, userID)
	if err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	// Load exercise information for each workout
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

	// Construct the response body
	response := struct {
		Session   storage.Session    `json:"session"`
		Exercises []storage.Exercise `json:"exercises"`
	}{
		Session:   *session,
		Exercises: exercises,
	}

	WriteJSON(w, http.StatusAccepted, response)
}

// HandleGetWorkouts Get all the workouts created by the user
func (a *Application) HandleGetWorkouts(w http.ResponseWriter, r *http.Request) {
	// get the userID. userID is parsed in the auth middleware
	userID := r.Context().Value("userID").(int)
	username := r.PathValue("user")

	// Database Call
	sessions, err := a.DB.Workout.GetWorkouts(userID)
	if err != nil && err != storage.ErrNoRecord {
		a.ServerError(w, err)
		return
	}

	// Construct Response
	res := &storage.WorkoutGetResponse{
		Username: username,
		Sessions: sessions,
	}

	WriteJSON(w, http.StatusOK, res)
}

// HandleUpdateWorkout Updates a specific workout session for a specific user
func (a *Application) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	// Initialize WorkoutRequest object to parse the request body
	var workoutRequest storage.WorkoutRequest

	// Parse the request body using json, returns error at failure
	if err := json.NewDecoder(r.Body).Decode(&workoutRequest); err != nil {
		a.BadRequest(w)
		return
	}

	defer r.Body.Close()

	// Validate the response
	if err := a.Validate.Struct(workoutRequest); err != nil {
		a.BadRequest(w)
		return
	}

	// get the userID. userID is parsed in the auth middleware
	userID := r.Context().Value("userID").(int)

	// get the sessionID from the path
	sessionID, err := strconv.Atoi(r.PathValue("workout"))
	if err != nil {
		a.BadRequest(w)
		return
	}

	// Database Call
	if err := a.DB.Workout.UpdateWorkout(userID, sessionID, workoutRequest); err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		}
		return
	}

	// Redirect to the newly updated workout
	http.Redirect(w, r, fmt.Sprintf("/users/%d/workouts/%d", userID, sessionID), http.StatusPermanentRedirect)
}

// HandleDeleteWorkout Deletes a specific workout
func (a *Application) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	// get the userID. userID is parsed in the auth middleware
	userID := r.Context().Value("userID").(int)

	// get the sessionID from the path
	sessionID, err := strconv.Atoi(r.PathValue("workout"))
	if err != nil {
		a.BadRequest(w)
		return
	}

	// Database Call
	if err := a.DB.Workout.DeleteWorkout(userID, sessionID); err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	// Construct response
	response := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("Workout session %d has been removed", sessionID),
	}

	WriteJSON(w, http.StatusAccepted, response)
}
