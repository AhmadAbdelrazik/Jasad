package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ahmadabdelrazik/jasad/pkg/validator"
)

type Workout struct {
	ID                int               `json:"id"`
	OwnerID           int               `json:"-"`
	Name              string            `json:"name"`
	Exercises         []WorkoutExercise `json:"exercises"`
	NumberOfExercises int               `json:"number_of_exercises,omitempty"`
	Version           int               `json:"-"`
}

func (w Workout) Validate(v *validator.Validator) {
	v.Check(w.Name != "", "name", "must not be empty")
	v.Check(len(w.Name) <= 50, "name", "must not be more than 50 bytes")
	v.Check(w.NumberOfExercises != 0, "exercises", "must include at least one exercise")
	if w.NumberOfExercises != len(w.Exercises) {
		panic("number of exercises does not match with the actual exercises")
	}

	for i, exercise := range w.Exercises {
		v.Check(i+1 == exercise.Order, "order", "exercises are not ordered correctly")
		exercise.Validate(v)
	}
}

type WorkoutExercise struct {
	ID       int       `json:"id"`
	Order    int       `json:"order"` // order in the workout
	Exercise *Exercise `json:"exercise"`
	Sets     int       `json:"sets"`
	Reps     int       `json:"reps,omitempty"`
	Weights  float32   `json:"weights,omitempty"`

	// in seconds
	RestAfter int  `json:"rest_after,omitempty"`
	Done      bool `json:"done"`
	Version   int  `json:"-"`
}

func (e WorkoutExercise) Validate(v *validator.Validator) {
	if e.Exercise == nil {
		panic("exercise not found")
	}
	v.Check(e.Sets > 0, "sets", "must be a positive number")
	v.Check(e.Sets < 1000, "sets", "must be less than 1000")

	v.Check(e.Reps >= 0, "reps", "must be a positive number")
	v.Check(e.Reps < 1000, "reps", "must be less than 1000")

	v.Check(e.Weights >= 0, "weights", "must be a positive number")
	v.Check(e.Weights < 1000, "weights", "must be less than 1000")

	v.Check(e.RestAfter >= 0, "rest_after", "must be a positive number")
	v.Check(e.RestAfter < 15*60, "rest_after", "must be less than 1000")
}

type WorkoutRepository struct {
	db *sql.DB
}

func (r *WorkoutRepository) Create(workout *Workout) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("can't start transaction: %w", err)
	}

	query := `
	INSERT INTO workouts(owner_id, name)
	VALUES ($1, $2)
	RETURNING id, version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = tx.QueryRowContext(ctx, query, workout.OwnerID, workout.Name).Scan(
		&workout.ID,
		&workout.Version,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting to workouts: %w", err)
	}

	query = `
	INSERT INTO workouts_exercises(workout_id, exercise_id, exercise_order,
	sets, reps, weights, rest_after, done)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, version
	`

	for i, exercise := range workout.Exercises {
		args := []any{
			workout.ID,
			exercise.Exercise.ID,
			exercise.Order,
			exercise.Sets,
			exercise.Reps,
			exercise.Weights,
			exercise.RestAfter,
			exercise.Done,
		}
		err := tx.QueryRowContext(ctx, query, args...).Scan(
			&workout.Exercises[i].ID,
			&workout.Exercises[i].Version,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting to workouts_exercises: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *WorkoutRepository) GetAllByID(ownerID int) ([]*Workout, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	// get basic info of each workout (not including exercises)
	query := `
	SELECT id, owner_id, name, version
	FROM workouts
	WHERE owner_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := tx.QueryContext(ctx, query, ownerID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	workouts := make([]*Workout, 0)

	for rows.Next() {
		workout := Workout{OwnerID: ownerID}
		err := rows.Scan(
			&workout.ID,
			&workout.OwnerID,
			&workout.Name,
			&workout.Version,
		)

		if err != nil {
			return nil, err
		}
		workouts = append(workouts, &workout)
	}

	if rows.Err() != nil {
		return nil, err
	}

	// get the exercises for each workout
	query = `
	SELECT we.id, we.exercise_order, we.sets, we.reps, we.weights,
	we.rest_after, we.done, we.version, e.id, e.name, e.muscle,
	e.instructions, e.additional_info, e.image_url, e.version
	FROM workouts_exercises AS we
	JOIN exercises AS e ON e.id = we.exercise_id
	WHERE we.workout_id = $1
	`

	for _, workout := range workouts {
		rows, err := tx.QueryContext(ctx, query, workout.ID)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		var exercises []WorkoutExercise

		for rows.Next() {
			var workoutExercise WorkoutExercise
			var exercise Exercise

			err := rows.Scan(
				&workoutExercise.ID,
				&workoutExercise.Order,
				&workoutExercise.Sets,
				&workoutExercise.Reps,
				&workoutExercise.Weights,
				&workoutExercise.RestAfter,
				&workoutExercise.Done,
				&workoutExercise.Version,
				&exercise.ID,
				&exercise.Name,
				&exercise.Muscle,
				&exercise.Instructions,
				&exercise.AdditionalInfo,
				&exercise.ImageURL,
				&exercise.Version,
			)

			if err != nil {
				return nil, err
			}

			workoutExercise.Exercise = &exercise
			exercises = append(exercises, workoutExercise)
		}

		workout.Exercises = exercises
		workout.NumberOfExercises = len(workout.Exercises)
	}

	return workouts, nil
}

func (r *WorkoutRepository) GetWorkoutByID(ownerID, workoutID int) (*Workout, error) {
	query := `
	SELECT w.name, w.version, we.id, we.exercise_order, we.sets,
	we.reps, we.weights, we.rest_after, we.done, we.version, e.id, e.name,
	e.muscle, e.instructions, e.additional_info, e.image_url, e.version
	FROM workouts AS w
	JOIN workouts_exercises AS we ON w.id = we.workout_id
	JOIN exercises AS e ON we.exercise_id = e.id
	WHERE w.owner_id = $1 AND w.id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, ownerID, workoutID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	workout := &Workout{
		OwnerID: ownerID,
		ID:      workoutID,
	}

	for rows.Next() {
		var workoutExercise WorkoutExercise
		var exercise Exercise

		err := rows.Scan(
			&workout.Name,
			&workout.Version,
			&workoutExercise.ID,
			&workoutExercise.Order,
			&workoutExercise.Sets,
			&workoutExercise.Reps,
			&workoutExercise.Weights,
			&workoutExercise.RestAfter,
			&workoutExercise.Done,
			&workoutExercise.Version,
			&exercise.ID,
			&exercise.Name,
			&exercise.Muscle,
			&exercise.Instructions,
			&exercise.AdditionalInfo,
			&exercise.ImageURL,
			&exercise.Version,
		)
		if err != nil {
			return nil, err
		}

		workoutExercise.Exercise = &exercise
		workout.Exercises = append(workout.Exercises, workoutExercise)
	}
	workout.NumberOfExercises = len(workout.Exercises)

	if workout.NumberOfExercises == 0 {
		return nil, ErrNotFound
	}

	return workout, nil
}

func (r *WorkoutRepository) Update(workout *Workout) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `
	UPDATE workouts SET name = $1, version = version + 1
	WHERE id = $2 AND version = $3
	RETURNING version
	`
	args := []any{workout.Name, workout.ID, workout.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = tx.QueryRowContext(ctx, query, args...).Scan(&workout.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	query = `DELETE from workouts_exercises WHERE workout_id = $1`
	_, err = tx.ExecContext(ctx, query, workout.ID)
	if err != nil {
		return err
	}

	query = `
	INSERT INTO workouts_exercises(workout_id, exercise_id, exercise_order, sets, reps, weights, rest_after, done)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, version
	`

	for i, exercise := range workout.Exercises {
		args := []any{
			workout.ID,
			exercise.Exercise.ID,
			exercise.Order,
			exercise.Sets,
			exercise.Reps,
			exercise.Weights,
			exercise.RestAfter,
			exercise.Done,
		}
		err := tx.QueryRowContext(ctx, query, args...).Scan(
			&workout.Exercises[i].ID,
			&workout.Exercises[i].Version,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func (r *WorkoutRepository) Delete(ownerID, workoutID int) error {
	query := `DELETE FROM workouts WHERE id = $1 AND owner_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, workoutID, ownerID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrNotFound
	}

	return nil
}
