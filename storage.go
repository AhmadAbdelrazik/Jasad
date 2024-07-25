package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var ErrNoRecord = errors.New("no records found")

type Storage interface {
	CreateExercise(*CreateExerciseRequest) error
	GetExercise(int) (*Exercise, error)
	UpdateExercise(*CreateExerciseRequest) error
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

	stmt = `INSERT INTO muscles_exercises(muscle_name, muscle_group, exercise_id), VALUES (?,?,?)`

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

func (st *MySQL) GetExercises(rows *sql.Rows) ([]Exercise, error) {
	return nil, nil
}

func (st *MySQL) GetExercise(ID int) (*Exercise, error) {
	return nil, nil
}
func (st *MySQL) UpdateExercise(ExerciseRequest *CreateExerciseRequest) error {
	return nil
}
func (st *MySQL) DeleteExercise(ID int) error {
	return nil
}
