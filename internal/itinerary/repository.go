package itinerary

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
	Create(itinerary *Itinerary) error
	GetByID(id int64) (*Itinerary, error)
	GetByTripID(tripID int64) ([]*Itinerary, error)
	Update(itinerary *Itinerary) error
	Delete(id int64) error
}

var _ RepositoryInterface = (*Repository)(nil)

func (r *Repository) Create(itinerary *Itinerary) error {
	query := `
        INSERT INTO itineraries (trip_id, title, description, date, created_by, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		itinerary.TripID,
		itinerary.Title,
		itinerary.Description,
		itinerary.Date,
		itinerary.CreatedBy,
		time.Now(),
		time.Now(),
	).Scan(&itinerary.ID)

	if err != nil {
		return fmt.Errorf("failed to create itinerary: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(id int64) (*Itinerary, error) {
	query := `
        SELECT id, trip_id, title, description, date, created_by, created_at, updated_at
        FROM itineraries
        WHERE id = $1`

	var itinerary Itinerary
	err := r.db.QueryRow(query, id).Scan(
		&itinerary.ID,
		&itinerary.TripID,
		&itinerary.Title,
		&itinerary.Description,
		&itinerary.Date,
		&itinerary.CreatedBy,
		&itinerary.CreatedAt,
		&itinerary.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("itinerary not found")
		}
		return nil, fmt.Errorf("failed to get itinerary: %w", err)
	}

	return &itinerary, nil
}

func (r *Repository) GetByTripID(tripID int64) ([]*Itinerary, error) {
	query := `
        SELECT id, trip_id, title, description, date, created_by, created_at, updated_at
        FROM itineraries
        WHERE trip_id = $1
        ORDER BY date ASC`

	rows, err := r.db.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get itineraries: %w", err)
	}
	defer rows.Close()

	var itineraries []*Itinerary
	for rows.Next() {
		var itinerary Itinerary
		err := rows.Scan(
			&itinerary.ID,
			&itinerary.TripID,
			&itinerary.Title,
			&itinerary.Description,
			&itinerary.Date,
			&itinerary.CreatedBy,
			&itinerary.CreatedAt,
			&itinerary.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan itinerary: %w", err)
		}
		itineraries = append(itineraries, &itinerary)
	}

	return itineraries, nil
}

func (r *Repository) Update(itinerary *Itinerary) error {
	query := `
        UPDATE itineraries
        SET title = $1, description = $2, date = $3, updated_at = $4
        WHERE id = $5`

	_, err := r.db.Exec(
		query,
		itinerary.Title,
		itinerary.Description,
		itinerary.Date,
		time.Now(),
		itinerary.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update itinerary: %w", err)
	}

	return nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM itineraries WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete itinerary: %w", err)
	}

	return nil
}
