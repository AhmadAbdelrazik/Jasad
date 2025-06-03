package model

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/ahmadabdelrazik/jasad/pkg/validator"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    Role   `json:"role"`
	Version int    `json:"-"`
}

func (u User) Validate(v *validator.Validator) {

}

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) Create(user *User) error {
	query := `
	INSERT INTO users(name, email, role)
	VALUES($1, $2, $3)
	RETURNING id, version
	`
	args := []any{user.Name, user.Email, user.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value"):
			return ErrAlreadyExists
		default:
			return err
		}
	}

	return nil
}

func (r *UserRepository) GetByID(id int) (*User, error) {
	query := `
	SELECT id, name, email, role, version
	FROM users
	WHERE id = $1
	`

	user := &User{
		ID: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, name, email, role, version
	FROM users
	WHERE email = $1
	`

	user := &User{
		Email: email,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (r *UserRepository) GetAll() ([]*User, error) {
	query := `
	SELECT id, name, email, role, version
	FROM users
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.Version,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
