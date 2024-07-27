package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var ErrNoRecord = errors.New("no records found")

type Storage interface {
	CreateExercise(*CreateExerciseRequest) error
	GetExercise(string) (*Exercise, error)
	GetExercises() ([]Exercise, error)
	UpdateExercise(*UpdateExerciseRequest) error
	DeleteExercise(int) error
	MuscleExists(*Muscle) error
}

type MySQL struct {
	DB *sql.DB
}

func NewMySQLServer() (*MySQL, error) {
	dsn := `ahmad:password@/jasad`
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &MySQL{DB: db}, nil
}

func (st *MySQL) MuscleExists(muscle *Muscle) error {
	tx, err := st.DB.Begin()
	if err != nil {
		return err
	}

	stmt := `SELECT muscle_name, muscle_group FROM muscles WHERE muscle_name = ? AND muscle_group = ?`
	row := tx.QueryRow(stmt, muscle.MuscleName, muscle.MuscleGroup)
	m := &Muscle{}
	if err := row.Scan(&m.MuscleName, &m.MuscleGroup); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}
	return nil
}

func (st *MySQL) CreateExercise(ExerciseRequest *CreateExerciseRequest) error {
	tx, err := st.DB.Begin()
	if err != nil {
		return err
	}

	for _, muscle := range ExerciseRequest.Muscles {
		if err := st.MuscleExists(&muscle); err != nil {
			return err
		}
	}

	stmt := `INSERT INTO exercises(exercise_name, exercise_description, reference_video) VALUES (?,?,?)`

	result, err := tx.Exec(stmt, ExerciseRequest.ExerciseName, ExerciseRequest.ExerciseDescription, ExerciseRequest.ReferenceVideo)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt = `INSERT INTO muscles_exercises(muscle_name, muscle_group, exercise_id) VALUES (?,?,?)`

	for _, muscle := range ExerciseRequest.Muscles {
		if _, err := tx.Exec(stmt, muscle.MuscleName, muscle.MuscleGroup, int(id)); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (st *MySQL) GetExercises() ([]Exercise, error) {
	tx, err := st.DB.Begin()
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

	exercises := []Exercise{}

	for rows.Next() {
		exercise := Exercise{}
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

	stmt = `SELECT muscle_name, muscle_group from muscles_exercises WHERE exercise_id = ?`

	for i := range exercises {
		rows, err := tx.Query(stmt, exercises[i].ExerciseID)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			} else {
				return nil, err
			}
		}

		var muscles []Muscle

		for rows.Next() {
			var muscle Muscle
			if err := rows.Scan(&muscle.MuscleName, &muscle.MuscleGroup); err != nil {
				tx.Rollback()
				return nil, err
			}
			muscles = append(muscles, muscle)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		exercises[i].Muscles = muscles
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return exercises, nil
}

func (st *MySQL) GetExercise(name string) (*Exercise, error) {
	return nil, nil
}
func (st *MySQL) UpdateExercise(Exercise *UpdateExerciseRequest) error {
	for _, muscle := range Exercise.Muscles {
		if err := st.MuscleExists(&muscle); err != nil {
			return err
		}
	}

	tx, err := st.DB.Begin()
	if err != nil {
		return err
	}

	stmt := `UPDATE exercises
	SET 
	exercise_name = ?,
	exercise_description = ?,
	reference_video = ?
	WHERE exercise_id = ?`

	_, err = tx.Exec(stmt, Exercise.ExerciseName, Exercise.ExerciseDescription, Exercise.ReferenceVideo, Exercise.ExerciseID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}

	stmt = `DELETE FROM muscles_exercises WHERE exercise_id = ?`
	_, err = tx.Exec(stmt, Exercise.ExerciseID)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt = `INSERT INTO muscles_exercises(exercise_id, muscle_name, muscle_group) VALUES (?,?,?)`
	for _, muscle := range Exercise.Muscles {
		if _, err := tx.Exec(stmt, Exercise.ExerciseID, muscle.MuscleName, muscle.MuscleGroup); err != nil {
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

func (st *MySQL) DeleteExercise(ID int) error {
	return nil
}
