package model

import (
	"database/sql"
	"errors"
)

type Muscle struct {
	MuscleName  string `json:"muscle_name"`
	MuscleGroup string `json:"muscle_group"`
}

type MuscleModel struct {
	DB *sql.DB
}

func (m *MuscleModel) GetAll() ([]Muscle, error) {
	stmt := `SELECT muscle_name, muscle_group FROM muscles`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	defer rows.Close()

	var muscles []Muscle

	for rows.Next() {

		var muscle Muscle
		err := rows.Scan(&muscle.MuscleName, &muscle.MuscleGroup)
		if err != nil {
			return nil, err
		}

		muscles = append(muscles, muscle)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return muscles, nil
}

func (m *MuscleModel) AllExist(muscleName []string) (string, error) {
	stmt := `SELECT muscle_name, muscle_group FROM muscles WHERE muscle_name = ?`

	for _, muscle := range muscleName {
		row := m.DB.QueryRow(stmt, muscle)

		var m Muscle

		err := row.Scan(&m.MuscleName, &m.MuscleGroup)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return muscle, ErrNoRecord
			} else {
				return "", err
			}
		}
	}
	return "", nil
}
