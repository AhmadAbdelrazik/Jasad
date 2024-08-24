package storage

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

// CreateUser takes the User request and creates a new user.
// returns error in case of already existing username or
// server error. Returns the UserJWT for token issuing.
func (st *MySQL) CreateUser(user *UserRequest) (*UserJWT, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return nil, err
	}

	hash, err := HashUserPasswords(user.Password, user.UserName)
	if err != nil {
		return nil, err
	}

	stmt := `INSERT INTO users(user_name, password, role, created_at) VALUES (?,?,?,?)`

	_, err = tx.Exec(stmt, user.UserName, hash, "user", time.Now())
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, ErrDuplicateEntry
		} else {
			return nil, err
		}
	}

	userJWT := &UserJWT{
		UserName: user.UserName,
		Role:     "user",
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return userJWT, nil
}

// CheckUserExists checks if the login credentials are valid.
// Returns invalid credentials if not found or server error.
// in case of success returns a UserJWT for token issuing.
func (st *MySQL) CheckUserExists(user *UserRequest) (*UserJWT, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return nil, err
	}

	stmt := `SELECT password, role FROM users WHERE user_name = ?`

	row := tx.QueryRow(stmt, user.UserName)

	var hash, role string

	if err := row.Scan(&hash, &role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidUsername
		} else {
			return nil, err
		}
	}

	check, err := CompareProvidedPassword(user.Password, user.UserName, hash)
	if err != nil {
		return nil, err
	}

	if !check {
		return nil, ErrInvalidPassword
	}

	return &UserJWT{UserName: user.UserName, Role: role}, nil
}

// GetUserID Gets UserID using the username. It's preferable to
// use userID for internal operations, and use username only
// for authentication. returns ErrNoRecord if the username is
// not linked with any User. or server error if happened.
func (st *MySQL) GetUserID(username string) (int, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return 0, err
	}

	stmt := `SELECT user_id FROM users WHERE user_name = ?`

	row := tx.QueryRow(stmt, username)

	var userId int

	if err := row.Scan(&userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}

	return userId, nil
}

func (st *MySQL) GetUser(userID int) (*UserResponse, error) {
	tx, err := st.DB.Begin()
	if err != nil {
		return nil, err
	}

	user := &UserResponse{
		UserID: userID,
	}
	stmt := `SELECT user_name, role, created_at FROM users WHERE user_id = ?`

	row := tx.QueryRow(stmt, userID)

	if err := row.Scan(stmt, &user.UserName, &user.Role, &user.CreatedAt); err != nil {
		return nil, err
	}

	return user, nil
}
