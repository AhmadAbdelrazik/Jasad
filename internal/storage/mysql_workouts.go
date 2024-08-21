package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (st *MySQL) CreateWorkout(workout WorkoutCreateRequest, userID int) (int, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO sessions(user_id, date) VALUES (?,?)`

	res, err := tx.Exec(stmt, userID, workout.Date.Round(24*time.Hour))
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sessionID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	ch := make(chan error, len(workout.Workouts))

	stmt = `INSERT INTO workouts(session_id, exercise_id, reps, sets, weights) VALUES (?, ?, ?, ?, ?)`
	for i, w := range workout.Workouts {
		fmt.Printf("w: %v\n", w)
		go func() {
			_, err := tx.Exec(stmt, sessionID, w.ExerciseID, w.Reps, w.Sets, w.Weights)
			if err != nil {
				ch <- fmt.Errorf("error at index %d: %w", i, err)
			}
			ch <- nil
		}()
	}

	for range workout.Workouts {
		if err := <-ch; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return int(sessionID), nil
}

func (st *MySQL) GetWorkout(sessionID, userID int) (*Session, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return nil, err
	}

	var session Session = Session{
		SessionID: sessionID,
		UserID:    userID,
	}

	stmt := `SELECT date FROM sessions WHERE user_id = ? AND session_id = ?`
	row := tx.QueryRow(stmt, userID, sessionID)

	if err := row.Scan(&session.Date); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	stmt = `SELECT workout_id, exercise_id, reps, sets, weights FROM workouts WHERE session_id = ?`

	rows, err := tx.Query(stmt, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	for rows.Next() {
		var w Workout
		if err := rows.Scan(&w.WorkoutID, &w.ExerciseID, &w.Reps, &w.Sets, &w.Weights); err != nil {
			return nil, err
		}

		session.Workouts = append(session.Workouts, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &session, nil
}

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
