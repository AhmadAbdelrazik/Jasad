package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Workout struct {
	ID                int               `json:"id"`
	UserID            int               `json:"-"`
	Name              string            `json:"name"`
	Exercises         []WorkoutExercise `json:"exercises"`
	NumberOfExercises int               `json:"number_of_exercises,omitempty"`
	Version           int               `json:"-"`
}

type WorkoutExercise struct {
	ID       int       `json:"id"`
	Order    int       `json:"order"` // order in the workout
	Exercise *Exercise `json:"exercise"`
	Sets     int       `json:"sets"`
	Reps     int       `json:"reps,omitempty"`
	Weights  int       `json:"weights,omitempty"`

	// in seconds
	RestAfter int  `json:"rest_after,omitempty"`
	Done      bool `json:"done"`
	Version   int  `json:"-"`
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
	INSERT INTO workouts(name)
	VALUES ($1)
	RETURNING id, version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = tx.QueryRowContext(ctx, query, workout.Name).Scan(
		&workout.ID,
		&workout.Exercises,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `
	INSERT INTO workouts_users(workout_id, user_id)
	VALUES ($1, $2)
	`

	result, err := tx.ExecContext(ctx, query, workout.ID, workout.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	n, err := result.RowsAffected()
	if err != nil || n == 0 {
		tx.Rollback()
		return err
	}

	query = `
	INSERT INTO workouts_exercises(workout_id, exercise_id, order, sets, reps, weights, rest_after, done)
	VALUES($1, $2, $3, $4, $5, $6, $7, &8)
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

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
