package storage

import (
	"database/sql"
	"errors"
)

// Checks if muscle Exists in the DB.
// Returns ErrNoRows if not found, or error if something went wrong.
// Returns nil if it's found
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
			return ErrInvalidMuscle
		} else {
			return err
		}
	}
	return nil
}
