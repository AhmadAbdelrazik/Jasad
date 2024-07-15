package model

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Exercise struct {
	ExerciseID          uuid.UUID `json:"exercise_id"`
	ExerciseName        string    `json:"exercise_name"`
	ExerciseDescription string    `json:"exercise_explanation"`
	ReferenceVideo      string    `json:"reference_video"`
	Muscles             []string  `json:"muscle_name"`
}

type ExerciseModel struct {
	DB *sql.DB
}

func (e *ExerciseModel) Create(exerciseName, exerciseDescription, ReferenceVideo string, muscles []string) (uuid.UUID, error) {
	var emptyUUID uuid.UUID
	id := uuid.New()

	tx, err := e.DB.Begin()

	if err != nil {
		return emptyUUID, err
	}
	stmt1 := `INSERT INTO exercises(exercise_id, exercise_name, exercise_description, reference_video)
VALUES(?, ?, ?, ?)`

	// Add the exercise itself

	_, err = tx.Exec(stmt1, id, exerciseName, exerciseDescription, ReferenceVideo)
	if err != nil {
		tx.Rollback()
		return emptyUUID, err
	}

	// Make the realtion between exercises and muscles
	stmt2 := `INSERT INTO muscles_exercises(exercise_id, muscle_name) VALUES(?, ?)`

	for _, muscle := range muscles {
		_, err := tx.Exec(stmt2, id, muscle)
		if err != nil {
			tx.Rollback()
			return emptyUUID, err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return emptyUUID, err
	}
	return id, nil
}

func (e *ExerciseModel) GetByID(id uuid.UUID) (*Exercise, error) {

	tx, err := e.DB.Begin()

	if err != nil {
		return nil, err
	}

	stmt := `SELECT exercise_name, exercise_description, reference_video
	FROM exercises WHERE exercise_id = ?`

	row := tx.QueryRow(stmt, id)

	exercise := Exercise{
		ExerciseID: id,
	}

	err = row.Scan(&exercise.ExerciseName, &exercise.ExerciseDescription, &exercise.ReferenceVideo)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	stmt = `SELECT muscle_name FROM muscles_exercises WHERE exercise_id = ?`

	rows, err := tx.Query(stmt, id)

	// TODO : ErrNoRecord Implementation here and in the handler.
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	for rows.Next() {
		var muscle string
		err := rows.Scan(&muscle)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		exercise.Muscles = append(exercise.Muscles, muscle)
	}
	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return &exercise, nil
}

func (e *ExerciseModel) GetAll() ([]Exercise, error) {
	tx, err := e.DB.Begin()
	if err != nil {
		return nil, err
	}

	stmt := `SELECT exercise_id, exercise_name, exercise_description, reference_video FROM exercises`

	rows, err := tx.Query(stmt)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	var exercises []Exercise

	for rows.Next() {
		var exercise Exercise
		err := rows.Scan(&exercise.ExerciseID, &exercise.ExerciseName, &exercise.ExerciseDescription, &exercise.ReferenceVideo)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		exercises = append(exercises, exercise)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	stmt = `SELECT muscle_name FROM muscles_exercises WHERE exercise_id = ?`

	for _, exercise := range exercises {
		rows, err = tx.Query(stmt, exercise.ExerciseID)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			} else {
				return nil, err
			}
		}

		for rows.Next() {
			var muscle string
			err := rows.Scan(&muscle)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			exercise.Muscles = append(exercise.Muscles, muscle)
		}

		if err := rows.Err(); err != nil {
			tx.Rollback()
			return nil, err
		}

	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return exercises, nil
}

func (e *ExerciseModel) Delete(id uuid.UUID) error {
	stmt := `DELETE FROM exercises WHERE id = ?`

	tx, err := e.DB.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec(stmt, id)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
			return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rows == 0 {
		return ErrNoRecord
	}

	return nil
}