package storage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (st *MySQL) CreateUser(user *UserCreateRequest) (string, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return "", err
	}

	hash, err := HashUserPasswords(user.Password, user.UserName)
	if err != nil {
		return "", err
	}

	userID := uuid.New()
	stmt := `INSERT INTO users(user_id, username, password, role, created_at) VALUES (?,?,?,?,?)`

	_, err = tx.Exec(stmt, userID.String(), user.UserName, hash, "user", time.Now())
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return userID.String(), nil
}

func (st *MySQL) CheckUserExists(user *UserSigninRequest) (string, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return "", err
	}

	stmt := `SELECT user_id, password FROM users WHERE username = ?`

	row := tx.QueryRow(stmt, user.UserName)

	var userID string
	var hash string

	if err := row.Scan(&userID, &hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoRecord
		} else {
			return "", err
		}
	}

	check, err := CompareProvidedPassword(user.Password, user.UserName, hash)
	if err != nil {
		return "", err
	}

	if check {
		return userID, nil
	} else {
		return "", ErrInvalidCredentials
	}
}

func (st *MySQL) GetUserByID(userID string) (*UserJWT, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return nil, err
	}

	stmt := `SELECT username, role FROM users WHERE user_id = ?`

	row := tx.QueryRow(stmt, userID)

	user := &UserJWT{
		UserID: userID,
	}

	if err := row.Scan(&user.UserName, &user.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}
