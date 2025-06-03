package model

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("resource not found")
	ErrEditConflict  = errors.New("update conflict")
)

type Model struct {
	Exercises *ExerciseRepository
	Users     *UserRepository
	Tokens    *TokenRepository
}

func New(dsn string) (*Model, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	redis, err := newRedisClient()
	if err != nil {
		return nil, err
	}

	return &Model{
		Exercises: &ExerciseRepository{db: db},
		Users:     &UserRepository{db: db},
		Tokens:    &TokenRepository{redis: redis},
	}, nil
}
