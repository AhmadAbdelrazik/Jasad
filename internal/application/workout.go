package application

import (
	"errors"
	"net/http"

	"github.com/ahmadabdelrazik/jasad/internal/model"
	"github.com/ahmadabdelrazik/jasad/pkg/validator"
)

func (app *Application) createWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	// get user id
	user, ok := getUser(r)
	if !ok {
		UnauthorizedResponse(w, r)
		return
	}

	var input struct {
		Name      string `json:"name"`
		Exercises []struct {
			Order      int     `json:"order"` // order in the workout
			ExerciseID int     `json:"exercise_id"`
			Sets       int     `json:"sets"`
			Reps       int     `json:"reps,omitempty"`
			Weights    float32 `json:"weights,omitempty"`

			// in seconds
			RestAfter int  `json:"rest_after,omitempty"`
			Done      bool `json:"done"`
		} `json:"exercises"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	// get all the exercise IDs from input
	ids := make([]int, len(input.Exercises))
	for i := range ids {
		ids[i] = input.Exercises[i].ExerciseID
	}

	// validate number of exercises has not reached limit of 20 exercise
	if len(ids) > 20 {
		v := validator.New()
		v.AddError("exercises", "must be less than 20 exercise")
		FailedValidationResponse(w, r, v.Errors)
		return
	}

	// fetch all exercises by their IDs
	exercises, err := app.models.Exercises.GetByIDs(ids...)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			NotFoundResponse(w, r)
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	workoutExercises := make([]model.WorkoutExercise, len(input.Exercises))

	for i, we := range input.Exercises {
		workoutExercises[i].Order = we.Order
		workoutExercises[i].Sets = we.Sets
		workoutExercises[i].Reps = we.Reps
		workoutExercises[i].Weights = we.Weights
		workoutExercises[i].RestAfter = we.RestAfter
		workoutExercises[i].Done = we.Done
		workoutExercises[i].Exercise = exercises[i]
	}

	workout := &model.Workout{
		UserID:            user.ID,
		Name:              input.Name,
		Exercises:         workoutExercises,
		NumberOfExercises: len(workoutExercises),
	}

	v := validator.New()
	workout.Validate(v)

	if !v.Valid() {
		FailedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.Workouts.Create(workout); err != nil {
		ServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"workout": workout}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func (app *Application) getWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	// get user id
	user, ok := getUser(r)
	if !ok {
		UnauthorizedResponse(w, r)
		return
	}

	workouts, err := app.models.Workouts.GetAllByID(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			NotFoundResponse(w, r)
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	v := validator.New()
	for _, workout := range workouts {
		workout.Validate(v)
		if !v.Valid() {
			FailedValidationResponse(w, r, v.Errors)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workouts": workouts}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func (app *Application) getWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	// get user id
	user, ok := getUser(r)
	if !ok {
		UnauthorizedResponse(w, r)
		return
	}

	id, err := app.readIDParam(r)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	workout, err := app.models.Workouts.GetWorkoutByID(user.ID, int(id))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			NotFoundResponse(w, r)
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	v := validator.New()
	workout.Validate(v)
	if !v.Valid() {
		FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"workout": workout}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}
