package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

var ErrNoRecord = errors.New("no records found")

type Storage interface {
	// Create Operations
	CreateExercise(*CreateExerciseRequest) error
	// Get Operations
	GetExercises() ([]Exercise, error)
	GetExercisesByMuscle(Muscle) ([]Exercise, error)
	GetExerciseByID(int) (*Exercise, error)
	GetExerciseByName(string) (*Exercise, error)
	// Update Operations
	UpdateExercise(*UpdateExerciseRequest) error
	// Delete Operations
	DeleteExercise(int) error
	// Helpers
	MuscleExists(*Muscle) error
}

type MySQL struct {
	DB *sql.DB
}

// Initalize New MySQL Database, inject it in the APIServer instance.
// returns the MySQL Database or an error
func NewMySQLDatabase() (*MySQL, error) {
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

// Creates an exercise and store it in the database
// It accepts a CreateExerciseRequest struct
// returns nil in success or an error
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

// Get All Exercises
// return all exercises or an Error
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

// Get Exercises with specific Muscle
// Takes a muscle as an argument
// returns the exercises or an error
func (st *MySQL) GetExercisesByMuscle(muscle Muscle) ([]Exercise, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return nil, err
	}

	stmt := `SELECT exercise_id FROM muscles_exercises WHERE muscle_name = ? AND muscle_group = ?`

	rows, err := tx.Query(stmt, muscle.MuscleName, muscle.MuscleGroup)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	var IDs []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		IDs = append(IDs, id)
	}
	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	var Exercises []Exercise
	for _, id := range IDs {
		exercise, err := st.GetExerciseByID(id)
		if err != nil {
			return nil, err
		}

		Exercises = append(Exercises, *exercise)
	}

	return Exercises, nil
}

// Gets Exercise info by Exercise Name.
// Returns the Exercise if found or ErrNoRecord if not found.
// Returns error if something went wrong
func (st *MySQL) GetExerciseByName(exerciseName string) (*Exercise, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return nil, err
	}

	exercise := Exercise{ExerciseName: exerciseName}

	stmt := `SELECT exercise_id, exercise_description, reference_video FROM exercises WHERE exercise_name = ?`
	row := tx.QueryRow(stmt, exercise.ExerciseName)

	if err := row.Scan(&exercise.ExerciseID, &exercise.ExerciseDescription, &exercise.ReferenceVideo); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	stmt = `SELECT muscle_name, muscle_group FROM muscles_exercises WHERE exercise_id = ?`
	rows, err := tx.Query(stmt, exercise.ExerciseID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for rows.Next() {
		var muscle Muscle

		if err := rows.Scan(muscle.MuscleName, muscle.MuscleGroup); err != nil {
			tx.Rollback()
			return nil, err
		}

		exercise.Muscles = append(exercise.Muscles, muscle)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &exercise, nil
}

// Gets Exercise info by Exercise ID.
// Returns the Exercise if found or ErrNoRecord if not found.
// Returns error if something went wrong
func (st *MySQL) GetExerciseByID(ID int) (*Exercise, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return nil, err
	}

	exercise := Exercise{ExerciseID: ID}

	stmt := `SELECT exercise_name, exercise_description, reference_video FROM exercises WHERE exercise_id = ?`
	row := tx.QueryRow(stmt, exercise.ExerciseID)

	if err := row.Scan(&exercise.ExerciseName, &exercise.ExerciseDescription, &exercise.ReferenceVideo); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	stmt = `SELECT muscle_name, muscle_group from muscles_exercises WHERE exercise_id = ?`
	rows, err := tx.Query(stmt, exercise.ExerciseID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for rows.Next() {
		var muscle Muscle

		if err := rows.Scan(muscle.MuscleName, muscle.MuscleGroup); err != nil {
			tx.Rollback()
			return nil, err
		}

		exercise.Muscles = append(exercise.Muscles, muscle)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &exercise, nil
}

// Updates an Exercise using Exercise ID.
// Accepts an UpdateExerciseRequest struct.
// Returns nil upon success or error.
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

// Deletes an Exercise By ID.
// return nil on success, or ErrNoRecord if Exercise was not found.
// return error if something went wrong
func (st *MySQL) DeleteExercise(ID int) error {
	tx, err := st.DB.Begin()

	if err != nil {
		return err
	}

	stmt := `DELETE FROM muscles_exercises WHERE exercise_id = ?`
	if _, err := tx.Exec(stmt, ID); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}

	stmt = `DELETE FROM exercises WHERE exercise_id = ?`
	if _, err := tx.Exec(stmt, ID); err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Checks if muscle Exists in the DB.
// Returns nil if it's found
// Returns ErrNoRows if not found, or error if something went wrong.
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
