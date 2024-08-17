package storage

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var ErrNoRecord = errors.New("no records found")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrDuplicateEntry = errors.New("duplicate entry")
var ErrInvalidMuscle = errors.New("invalid muscle")

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
	// CreateExercise Add a new Exercise to the database.
	//
	// returns ErrInvalidMuscle if there was any invalid
	// muscles in the ExerciseCreateRequest. returns
	// server error on other failure cases, return nil at
	// success
	CreateExercise(*ExerciseCreateRequest) error

	// GetExercises Get all exercises in the database.
	//
	// Returns ErrNoRecord if the database is empty.
	// returns server error on other failure cases,
	// returns array of exercises.
	GetExercises() ([]Exercise, error)

	// GetExercisesByMuscle Returns all Exercises that
	// Contains Muscle in them.
	//
	// Returns ErrNoRecord if there is no muscles
	// returns server error on other failure cases,
	// returns array of exercises.
	GetExercisesByMuscle(Muscle) ([]Exercise, error)

	// GetExerciseByID Returns the Exercise with the
	// specific id.
	//
	// Returns ErrNoRecord if exercise with id does
	// not exist. returns server error on other failure
	// cases, returns the exercise.
	GetExerciseByID(int) (*Exercise, error)

	// GetExerciseByName Returns the Exercise with
	// the specific name.
	//
	// Returns ErrNoRecord if exercise with name does
	// not exist. returns server error on other failure
	// cases, returns the exercise.
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
