package storage

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var ErrNoRecord = errors.New("no records found")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrDuplicateEntry = errors.New("duplicate entry")

type Storage struct {
	User     IUserStorage
	Exercise IExerciseStorage
}

type IUserStorage interface {
	// CreateUser takes the User request and creates a new user.
	// returns ErrDuplicateEntry in case of already existing username
	// or returns server error. Returns the UserJWT for token issuing.
	CreateUser(*UserRequest) (*UserJWT, error)

	// CheckUserExists checks if the login credentials are valid.
	// Returns ErrInvalidCredentials if not found or server error.
	// in case of success returns a UserJWT for token issuing.
	CheckUserExists(*UserRequest) (*UserJWT, error)

	// GetUserID Gets UserID using the username. It's preferable to
	// use userID for internal operations, and use username only
	// for authentication. returns ErrNoRecord if the username is
	// not linked with any User. or server error if happened.
	GetUserID(username string) (int, error)

	// GetUser Gets All information about the user. returns
	// ErrNoRecord if user is not found, or error if server error.
	GetUser(userID int) (*User, error)
}

type IExerciseStorage interface {
	CreateExercise(*ExerciseCreateRequest) error

	GetExercises() ([]Exercise, error)
	GetExercisesByMuscle(Muscle) ([]Exercise, error)
	GetExerciseByID(int) (*Exercise, error)
	GetExerciseByName(string) (*Exercise, error)

	UpdateExercise(*ExerciseUpdateRequest) error

	DeleteExercise(int) error

	MuscleExists(*Muscle) error
}

type MySQL struct {
	DB *sql.DB
}

// Initalize New MySQL Database, inject it in the APIServer instance.
// returns the MySQL Database or an error
func NewMySQLDatabase(dsn string) (*Storage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	mysql := &MySQL{DB: db}

	storage := &Storage{
		User:     mysql,
		Exercise: mysql,
	}

	return storage, nil
}
