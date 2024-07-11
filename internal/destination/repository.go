package destination

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
	Create(destination *Destination) error
	GetByTripID(tripID int64) (*Destination, error)
	Update(destination *Destination) error
	Delete(tripID int64) error
}

var _ RepositoryInterface = (*Repository)(nil)

func (r *Repository) Create(destination *Destination) error {
	query := `
        INSERT INTO destinations (trip_id, name, country, city, description, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		destination.TripID,
		destination.Name,
		destination.Country,
		destination.City,
		destination.Description,
		time.Now(),
		time.Now(),
	).Scan(&destination.ID)

	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}

	return nil
}

func (r *Repository) GetByTripID(tripID int64) (*Destination, error) {
	query := `
        SELECT id, trip_id, name, country, city, description, created_at, updated_at
        FROM destinations
        WHERE trip_id = $1`

	var destination Destination
	err := r.db.QueryRow(query, tripID).Scan(
		&destination.ID,
		&destination.TripID,
		&destination.Name,
		&destination.Country,
		&destination.City,
		&destination.Description,
		&destination.CreatedAt,
		&destination.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("destination not found")
		}
		return nil, fmt.Errorf("failed to get destination: %w", err)
	}

	return &destination, nil
}

func (r *Repository) Update(destination *Destination) error {
	query := `
        UPDATE destinations
        SET name = $1, country = $2, city = $3, description = $4, updated_at = $5
        WHERE trip_id = $6`

	_, err := r.db.Exec(
		query,
		destination.Name,
		destination.Country,
		destination.City,
		destination.Description,
		time.Now(),
		destination.TripID,
	)

	if err != nil {
		return fmt.Errorf("failed to update destination: %w", err)
	}

	return nil
}

func (r *Repository) Delete(tripID int64) error {
	query := `DELETE FROM destinations WHERE trip_id = $1`

	_, err := r.db.Exec(query, tripID)
	if err != nil {
		return fmt.Errorf("failed to delete destination: %w", err)
	}

	return nil
}
