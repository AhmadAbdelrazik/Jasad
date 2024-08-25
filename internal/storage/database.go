package storage

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var ErrNoRecord = errors.New("no records found")
var ErrInvalidUsername = errors.New("invalid username")
var ErrInvalidPassword = errors.New("invalid password")
var ErrDuplicateEntry = errors.New("duplicate entry")
var ErrInvalidMuscle = errors.New("invalid muscle")

type Storage struct {
	User     IUserStorage
	Exercise IExerciseStorage
	Workout  IWorkoutStorage
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
	GetUser(userID int) (*UserResponse, error)
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

	// UpdateExercise Updates the exercise identified by
	// exerciseID.
	//
	// Returns ErrInvalidMuscle if any one of the muscles are invalid
	// Returns ErrNoRecord if exercise with exerciseID doesn't exist
	UpdateExercise(exerciseID int, exercise *ExerciseUpdateRequest) error

	// DeleteExercise Deletes the exercise identified by exerciseID
	//
	// Returns ErrNoRecord if the exercise doesn't exist.
	DeleteExercise(int) error

	// MuscleExists Check if muscle Exists.
	//
	// Returns ErrInvalidMuscle if muscle doesn't exist.
	MuscleExists(*Muscle) error
}

type IWorkoutStorage interface {
	// CreateWorkout takes a WorkoutCreateRequest and userID to add
	// a workout to the database. returns the sessionID at success,
	// returns 0 and err at failure.
	CreateWorkout(workout WorkoutRequest, userID int) (int, error)

	// GetWorkout takes a sessionID and userID, to get a specific
	// workout. the userID is required to prevent Broken Object
	// Level Authorization or BOLA. Function returns session in
	// case of success, or returns a ErrNoRecord if no records are
	// found. returns an err if there was other database related errors.
	GetWorkout(sessionID, userID int) (*Session, error)

	// GetWorkouts gets workout related to the userID, in case of
	// failure, returns ErrNoRecord in case of no records, or err
	// if there was a database error
	GetWorkouts(userID int) ([]SessionResponse, error)

	// UpdateWorkout Updates the workouts with the specified userID
	// and sessionID, returns ErrNoRecord if no session is found,
	// and returns nil at success. otherwise returns a generic error.
	UpdateWorkout(userID, sessionID int, workout WorkoutRequest) error

	// DeleteWorkout Delete the workout session with the specified
	// and sessionID, returns ErrNoRecord if no session is found,
	// and returns nil at success. otherwise returns a generic error.
	DeleteWorkout(userID, sessionID int) error
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
		Workout:  mysql,
	}

	return storage, nil
}
