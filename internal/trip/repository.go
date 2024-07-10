package trip

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type RepositoryInterface interface {
	Create(trip *Trip) error
	GetByID(id int64) (*Trip, error)
	Update(trip *Trip) error
	Delete(id int64) error
}

var _ RepositoryInterface = (*Repository)(nil)

func (r *Repository) Create(trip *Trip) error {
	query := `
		INSERT INTO trips (name, description, start_date, end_date, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := r.db.QueryRow(
		query,
		trip.Name,
		trip.Description,
		trip.StartDate,
		trip.EndDate,
		trip.CreatedBy,
		time.Now(),
		time.Now(),
	).Scan(&trip.ID)

	if err != nil {
		return fmt.Errorf("failed to create trip: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(id int64) (*Trip, error) {
	query := `
		SELECT id, name, description, start_date, end_date, created_by, created_at, updated_at
		FROM trips
		WHERE id = $1`

	var trip Trip
	err := r.db.QueryRow(query, id).Scan(
		&trip.ID,
		&trip.Name,
		&trip.Description,
		&trip.StartDate,
		&trip.EndDate,
		&trip.CreatedBy,
		&trip.CreatedAt,
		&trip.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trip not found")
		}
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	return &trip, nil
}

func (r *Repository) Update(trip *Trip) error {
	query := `
		UPDATE trips
		SET name = $1, description = $2, start_date = $3, end_date = $4, updated_at = $5
		WHERE id = $6`

	_, err := r.db.Exec(
		query,
		trip.Name,
		trip.Description,
		trip.StartDate,
		trip.EndDate,
		time.Now(),
		trip.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update trip: %w", err)
	}

	return nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM trips WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trip: %w", err)
	}

	return nil
}
