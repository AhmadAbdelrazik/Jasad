package model

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Model struct {
	Exercises *ExerciseRepository
}

func New(dsn string) (*Model, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Model{
		Exercises: &ExerciseRepository{db: db},
	}, nil
}
