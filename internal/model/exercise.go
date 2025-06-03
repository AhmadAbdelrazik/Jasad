package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ahmadabdelrazik/jasad/pkg/validator"
)

type Exercise struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Muscle Muscle `json:"muscle"`

	// Step by step brief instructions for the exercise
	Instructions string `json:"instructions"`

	// In Depth details about the exercise
	AdditionalInfo string `json:"additional_info"`

	ImageURL string `json:"image_url"`

	// Version is used for Version control in databae
	Version int `json:"-"`
}

func (e Exercise) Validate(v *validator.Validator) {
	v.Check(strings.Trim(e.Name, " ") != "", "name", "can't be empty")
	v.Check(len(e.Name) < 50, "name", "must be less than 50 bytes")

	v.Check(strings.Trim(e.Instructions, " ") != "", "instructions", "can't be empty")
	v.Check(len(e.Instructions) < 1000, "instructions", "must be less than 1000 bytes")

	v.Check(strings.Trim(e.AdditionalInfo, " ") != "", "additional_info", "can't be empty")
	v.Check(len(e.AdditionalInfo) < 10000, "additional_info", "must be less than 10000 bytes")

	v.Check(validator.URLRX.MatchString(e.ImageURL), "image_url", "must be a valid url")
}

type ExerciseRepository struct {
	db *sql.DB
}

func (r *ExerciseRepository) Create(exercise *Exercise) error {
	query := `
	INSERT INTO exercises(name, muscle, instructions, additional_info, image_url)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, version
	`
	args := []any{
		exercise.Name,
		exercise.Muscle,
		exercise.Instructions,
		exercise.AdditionalInfo,
		exercise.ImageURL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exercise.ID, &exercise.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate"):
			return ErrAlreadyExists
		default:
			return err
		}
	}

	return nil
}

func (r *ExerciseRepository) Get(id int) (*Exercise, error) {
	query := `
	SELECT id, name, muscle, instructions, additional_info, image_url, version
	FROM exercises
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exercise := &Exercise{ID: id}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Muscle,
		&exercise.Instructions,
		&exercise.AdditionalInfo,
		&exercise.ImageURL,
		&exercise.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return exercise, nil
}

func (r *ExerciseRepository) Search(name, muscle string, filters Filters) ([]*Exercise, Metadata, error) {
	// We Use COUNT(*) OVER() to get the total number for metadata. we
	// utilize postgres text search using to_tsvector for better string
	// search. limi and offset are calculated based on the page and page
	// size queries from the coming request.
	query := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, name, muscle, instructions, additional_info, image_url, version
	FROM exercises
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND (to_tsvector('simple', muscle) @@ plainto_tsquery('simple', $2) OR $2 = '')
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{name, muscle, filters.limit(), filters.offset()}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	exercises := []*Exercise{}

	for rows.Next() {
		var exercise Exercise

		err := rows.Scan(
			&totalRecords,
			&exercise.ID,
			&exercise.Name,
			&exercise.Muscle,
			&exercise.Instructions,
			&exercise.AdditionalInfo,
			&exercise.ImageURL,
			&exercise.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		exercises = append(exercises, &exercise)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)

	return exercises, metadata, nil
}

func (r *ExerciseRepository) Update(exercise *Exercise) error {
	query := `
	UPDATE exercises
	SET name = $1, muscle = $2, instructions = $3, additional_info = $4, image_url = $5, version = version + 1
	WHERE id = $6 AND version = $7
	RETURNING version
	`

	args := []any{
		exercise.Name,
		exercise.Muscle,
		exercise.Instructions,
		exercise.AdditionalInfo,
		exercise.ImageURL,
		exercise.ID,
		exercise.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exercise.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		case strings.Contains(err.Error(), "duplicate key value"):
			return ErrAlreadyExists
		default:
			return err
		}
	}

	return nil
}

func (r *ExerciseRepository) Delete(id int) error {
	query := `DELETE FROM exercises WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrNotFound
	}

	return nil
}
