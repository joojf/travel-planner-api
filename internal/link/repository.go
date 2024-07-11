package link

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
	Create(link *Link) error
	GetByTripID(tripID int64) ([]*Link, error)
	GetByID(id int64) (*Link, error)
	Update(link *Link) error
	Delete(id int64) error
}

var _ RepositoryInterface = (*Repository)(nil)

func (r *Repository) Create(link *Link) error {
	query := `
        INSERT INTO links (trip_id, title, url, description, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		link.TripID,
		link.Title,
		link.URL,
		link.Description,
		time.Now(),
		time.Now(),
	).Scan(&link.ID)

	if err != nil {
		return fmt.Errorf("failed to create link: %w", err)
	}

	return nil
}

func (r *Repository) GetByTripID(tripID int64) ([]*Link, error) {
	query := `
        SELECT id, trip_id, title, url, description, created_at, updated_at
        FROM links
        WHERE trip_id = $1`

	rows, err := r.db.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get links: %w", err)
	}
	defer rows.Close()

	var links []*Link
	for rows.Next() {
		var link Link
		err := rows.Scan(
			&link.ID,
			&link.TripID,
			&link.Title,
			&link.URL,
			&link.Description,
			&link.CreatedAt,
			&link.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}
		links = append(links, &link)
	}

	return links, nil
}

func (r *Repository) GetByID(id int64) (*Link, error) {
	query := `
        SELECT id, trip_id, title, url, description, created_at, updated_at
        FROM links
        WHERE id = $1`

	var link Link
	err := r.db.QueryRow(query, id).Scan(
		&link.ID,
		&link.TripID,
		&link.Title,
		&link.URL,
		&link.Description,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("link not found")
		}
		return nil, fmt.Errorf("failed to get link: %w", err)
	}

	return &link, nil
}

func (r *Repository) Update(link *Link) error {
	query := `
        UPDATE links
        SET title = $1, url = $2, description = $3, updated_at = $4
        WHERE id = $5`

	_, err := r.db.Exec(
		query,
		link.Title,
		link.URL,
		link.Description,
		time.Now(),
		link.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update link: %w", err)
	}

	return nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM links WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}

	return nil
}
