package storage

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var ErrNoRecord = errors.New("no records found")
var ErrInvalidCredentials = errors.New("no records found")

type Storage interface {
	// Create Operations
	CreateExercise(*ExerciseCreateRequest) error
	// Get Operations
	GetExercises() ([]Exercise, error)
	GetExercisesByMuscle(Muscle) ([]Exercise, error)
	GetExerciseByID(int) (*Exercise, error)
	GetExerciseByName(string) (*Exercise, error)
	// Update Operations
	UpdateExercise(*ExerciseUpdateRequest) error
	// Delete Operations
	DeleteExercise(int) error
	// Helpers
	MuscleExists(*Muscle) error

	CreateUser(*UserCreateRequest) (uuid.UUID, error)
	CheckUserExists(user *UserSigninRequest) (uuid.UUID, error)
}

type MySQL struct {
	DB *sql.DB
}

// Initalize New MySQL Database, inject it in the APIServer instance.
// returns the MySQL Database or an error
func NewMySQLDatabase(dsn string) (*MySQL, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &MySQL{DB: db}, nil
}
