package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Creates an exercise and store it in the database
// It accepts a ExerciseCreateRequest struct
// returns ErrInvalidMuscles if any muscles are invalid
// returns ErrDuplicateEntry if exercise name already exists
func (st *MySQL) CreateExercise(ExerciseRequest *ExerciseCreateRequest) error {
	tx, err := st.DB.Begin()
	if err != nil {
		return err
	}

	// Validate Muscles
	for _, muscle := range ExerciseRequest.Muscles {
		if err := st.MuscleExists(&muscle); err != nil {
			return ErrInvalidMuscle
		}
	}

	stmt := `INSERT INTO exercises(exercise_name, exercise_description, reference_video) VALUES (?,?,?)`

	result, err := tx.Exec(stmt, ExerciseRequest.ExerciseName, ExerciseRequest.ExerciseDescription, ExerciseRequest.ReferenceVideo)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "Duplicate entry") {
			return ErrDuplicateEntry
		} else {
			return err
		}
	}

	// Get the ExerciseID for muscles_exercises insertion
	exerciseID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt = `INSERT INTO muscles_exercises(muscle_name, muscle_group, exercise_id) VALUES (?,?,?)`

	// Utilize Go routines for muscle insertion to database. using channels to return errors if there is any
	ch := make(chan error, len(ExerciseRequest.Muscles))

	for i, m := range ExerciseRequest.Muscles {
		go func() {
			if _, err := tx.Exec(stmt, m.MuscleName, m.MuscleGroup, int(exerciseID)); err != nil {
				ch <- fmt.Errorf("error at index %d: %w", i, err)
			}
			ch <- nil
		}()
	}

	// Check for any received errors in the channel
	for range ExerciseRequest.Muscles {
		if err := <-ch; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit changes
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

	// Get all exercises
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

	// enumerate over the result to populate the exercises.
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

	ch := make(chan error, len(exercises)) // make error channels with the length of exercises array

	for i := range exercises {
		go func() {
			// get all muscles related to the exercise
			rows, err := tx.Query(stmt, exercises[i].ExerciseID)
			if err != nil {
				ch <- err
			}

			var muscles []Muscle

			// enumerate over them and add them to muscles slice
			for rows.Next() {
				var muscle Muscle
				if err := rows.Scan(&muscle.MuscleName, &muscle.MuscleGroup); err != nil {
					ch <- err
				}
				muscles = append(muscles, muscle)
			}

			// check for enumeration errors
			if err := rows.Err(); err != nil {
				ch <- err
			}

			// add the muscles to the related exercises
			exercises[i].Muscles = muscles

			ch <- nil
		}()
	}

	// range over the retrieved errors from the channel.
	for range exercises {
		if err := <-ch; err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			} else {
				return nil, err
			}
		}
	}

	// Commit changes
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

		if err := rows.Scan(&muscle.MuscleName, &muscle.MuscleGroup); err != nil {
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

		if err := rows.Scan(&muscle.MuscleName, &muscle.MuscleGroup); err != nil {
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
func (st *MySQL) UpdateExercise(exerciseID int, exercise *ExerciseUpdateRequest) error {
	for _, muscle := range exercise.Muscles {
		if err := st.MuscleExists(&muscle); err != nil {
			return ErrInvalidMuscle
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

	_, err = tx.Exec(stmt, exercise.ExerciseName, exercise.ExerciseDescription, exercise.ReferenceVideo, exerciseID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}

	stmt = `DELETE FROM muscles_exercises WHERE exercise_id = ?`
	_, err = tx.Exec(stmt, exerciseID)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt = `INSERT INTO muscles_exercises(exercise_id, muscle_name, muscle_group) VALUES (?,?,?)`
	for _, muscle := range exercise.Muscles {
		if _, err := tx.Exec(stmt, exerciseID, muscle.MuscleName, muscle.MuscleGroup); err != nil {
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
