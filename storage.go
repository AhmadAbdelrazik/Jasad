package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Storage interface {
	CreateExercise(ExerciseName, MuscleName, MuscleGroup, Description, ReferenceVideo string) error
	GetExercise(int) (*Exercise, error)
	UpdateExercise(ExerciseName, MuscleName, MuscleGroup, Description, ReferenceVideo string) error
	DeleteExercise(int) error
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

func GetExercises(rows *sql.Rows) ([]Exercise, error) {

	return nil, nil
}

func (m *MySQL) CreateExercise(ExerciseName, MuscleName, MuscleGroup, Description, ReferenceVideo string) error {
	return nil
}

func (m *MySQL) GetExercise(ID int) (*Exercise, error) {
	return nil, nil
}
func (m *MySQL) UpdateExercise(ExerciseName, MuscleName, MuscleGroup, Description, ReferenceVideo string) error {
	return nil
}
func (m *MySQL) DeleteExercise(ID int) error {
	return nil
}
