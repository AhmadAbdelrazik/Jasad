package model

import (
	"database/sql"
	"errors"
)

type JasadModel struct {
	DB *sql.DB
}

type Exercise struct {
	ExerciseID          int      `json:"exercise_id"`
	ExerciseName        string   `json:"exercise_name"`
	ExerciseExplanation string   `json:"exercise_explanation"`
	ReferenceVideo      string   `json:"reference_video"`
	Muscles             []string `json:"muscle_name"`
}

type Muscle struct {
	MuscleName  string `json:"muscle_name"`
	MuscleGroup string `json:"muscle_group"`
}

func (j *JasadModel) GetAllMuscles() ([]Muscle, error) {
	stmt := `SELECT muscle_name, muscle_group FROM muscles`

	rows, err := j.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	defer rows.Close()

	var muscles []Muscle

	for rows.Next() {

		var muscle Muscle
		err := rows.Scan(&muscle.MuscleName, &muscle.MuscleGroup)
		if err != nil {
			return nil, err
		}

		muscles = append(muscles, muscle)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return muscles, nil
}

func (j *JasadModel) CheckMusclesExist(muscleName []string) (string, error) {
	stmt := `SELECT muscle_name, muscle_group FROM muscles WHERE muscle_name = ?`

	for _, muscle := range muscleName {
		row := j.DB.QueryRow(stmt, muscle)

		var m Muscle

		err := row.Scan(&m.MuscleName, &m.MuscleGroup)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return muscle, ErrNoRecord
			} else {
				return "", err
			}
		}
	}
	return "", nil
}

func (j *JasadModel) AddExercise(exerciseName, exerciseDescription, ReferenceVideo string, muscles []string) (int, error) {

	tx, err := j.DB.Begin()

	if err != nil {
		return 0, err
	}
	stmt1 := `INSERT INTO exercises(exercise_name, exercise_description, reference_video)
VALUES(?, ?, ?)`

	// Add the exercise itself
	result, err := tx.Exec(stmt1, exerciseName, exerciseDescription, ReferenceVideo)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Make the realtion between exercises and muscles
	stmt2 := `INSERT INTO muscles_exercises(exercise_id, muscle_name) VALUES(?, ?)`

	for _, muscle := range muscles {
		_, err := tx.Exec(stmt2, id, muscle)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}
	return int(id), nil
}
