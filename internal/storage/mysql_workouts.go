package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// CreateWorkout Creates a new workout session and assign it to userID
func (st *MySQL) CreateWorkout(workout WorkoutRequest, userID int) (int, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return 0, err
	}

	// Create a new session.
	stmt := `INSERT INTO sessions(user_id, date) VALUES (?,?)`

	res, err := tx.Exec(stmt, userID, workout.Date.Round(24*time.Hour))
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get the sessionID
	sessionID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Utilizing go routines to Create new workouts and assign it to the session.
	ch := make(chan error, len(workout.Workouts))

	stmt = `INSERT INTO workouts(session_id, exercise_id, reps, sets, weights) VALUES (?, ?, ?, ?, ?)`
	for i, w := range workout.Workouts {
		go func() {
			_, err := tx.Exec(stmt, sessionID, w.ExerciseID, w.Reps, w.Sets, w.Weights)
			if err != nil {
				ch <- fmt.Errorf("error at index %d: %w", i, err)
			}
			ch <- nil
		}()
	}

	// Check for any errors
	for range workout.Workouts {
		if err := <-ch; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Commit the changes
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(sessionID), nil
}

// GetWorkout Gets the workout session using sessionID for a specific user
// GetWorkout takes care of BOLA attacks since it gets sessions that are
// connected with the userID only
func (st *MySQL) GetWorkout(sessionID, userID int) (*Session, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return nil, err
	}

	var session Session = Session{
		SessionID: sessionID,
		UserID:    userID,
	}

	// get the session date. also used to know if session exists
	stmt := `SELECT date FROM sessions WHERE user_id = ? AND session_id = ?`
	row := tx.QueryRow(stmt, userID, sessionID)

	if err := row.Scan(&session.Date); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Get all workouts attached to the session
	stmt = `SELECT workout_id, exercise_id, reps, sets, weights FROM workouts WHERE session_id = ?`

	rows, err := tx.Query(stmt, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	// enumerate the workouts and add them to the session object one by one
	for rows.Next() {
		var w Workout

		err := rows.Scan(&w.WorkoutID, &w.ExerciseID, &w.Reps, &w.Sets, &w.Weights)
		if err != nil {
			return nil, err
		}

		session.Workouts = append(session.Workouts, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &session, nil
}

// GetWorkouts Get all the sessions created by the userID
func (st *MySQL) GetWorkouts(userID int) ([]SessionResponse, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return nil, err
	}

	var sessions []SessionResponse

	stmt := `SELECT session_id, date FROM sessions WHERE user_id = ?`

	rows, err := tx.Query(stmt, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		s := SessionResponse{}
		if err := rows.Scan(&s.SessionID, &s.Date); err != nil {
			return nil, err
		}

		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

// UpdateWorkout updates the workouts associated with sessionID of userID
func (st *MySQL) UpdateWorkout(userID, sessionID int, workout WorkoutRequest) error {
	tx, err := st.DB.Begin()
	if err != nil {
		return err
	}

	// Update the date
	stmt := `UPDATE sessions
	SET date = ?
	WHERE user_id = ? AND session_id = ?`

	_, err = tx.Exec(stmt, workout.Date, userID, sessionID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		}
	}

	// Deletes all old associated workouts
	stmt = `DELETE FROM workouts WHERE user_id = ? AND session_id = ?`

	_, err = tx.Exec(stmt, userID, sessionID)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Adds all the new associated workouts
	stmt = `INSERT INTO workouts(session_id, exercise_id, reps, sets, weights) VALUES(?, ?, ?, ?, ?)`

	// uses go routines to insert workouts
	ch := make(chan error, len(workout.Workouts))

	for i, w := range workout.Workouts {
		go func() {
			_, err := tx.Exec(stmt, sessionID, w.ExerciseID, w.Reps, w.Sets, w.Weights)
			if err != nil {
				ch <- fmt.Errorf("error at %d: %w", i, err)
			} else {
				ch <- nil
			}
		}()
	}

	// check for any errors
	for range workout.Workouts {
		if err := <-ch; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit changes
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// DeleteWorkout Deletes the session identified by sessionID and userID
func (st *MySQL) DeleteWorkout(userID, sessionID int) error {
	tx, err := st.DB.Begin()
	if err != nil {
		return err
	}

	// Delete the workouts associated with sessions
	stmt := `DELETE FROM workouts WHERE session_id = ?`
	_, err = tx.Exec(stmt, sessionID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// delete the session
	stmt = `DELETE FROM sessions WHERE user_id = ? AND session_id = ?`

	res, err := tx.Exec(stmt, userID, sessionID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if there was any deleted rows
	rows, err := res.RowsAffected()

	if err != nil {
		tx.Rollback()
		return err
	} else if rows == 0 {
		tx.Rollback()
		return ErrNoRecord
	}

	// Commit changes
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
