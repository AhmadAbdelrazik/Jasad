package storage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (st *MySQL) CreateUser(user *UserCreateRequest) (uuid.UUID, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	hash, err := HashUserPasswords(user.Password, user.UserName)
	if err != nil {
		return uuid.Nil, err
	}

	userID := uuid.New()
	stmt := `INSERT INTO users(user_id, username, password, role, created_at) VALUES (?,?,?,?,?)`

	_, err = tx.Exec(stmt, userID.String(), user.UserName, hash, "user", time.Now())
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (st *MySQL) CheckUserExists(user *UserSigninRequest) (uuid.UUID, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	stmt := `SELECT user_id, password FROM users WHERE username = ?`

	row := tx.QueryRow(stmt, user.UserName)

	var userID uuid.UUID
	var hash string

	if err := row.Scan(&userID, &hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ErrNoRecord
		} else {
			return uuid.Nil, err
		}
	}

	check, err := CompareProvidedPassword(user.Password, user.UserName, hash)
	if err != nil {
		return uuid.Nil, err
	}

	if check {
		return userID, nil
	} else {
		return uuid.Nil, ErrInvalidCredentials
	}
}
