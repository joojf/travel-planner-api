package invitation

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
	Create(invitation *Invitation) error
	GetByTripID(tripID int64) ([]*Invitation, error)
	Delete(id int64) error
}

var _ RepositoryInterface = (*Repository)(nil)

func (r *Repository) Create(invitation *Invitation) error {
	query := `
        INSERT INTO invitations (trip_id, email, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		invitation.TripID,
		invitation.Email,
		invitation.Status,
		time.Now(),
		time.Now(),
	).Scan(&invitation.ID)

	if err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}

	return nil
}

func (r *Repository) GetByTripID(tripID int64) ([]*Invitation, error) {
	query := `
        SELECT id, trip_id, email, status, created_at, updated_at
        FROM invitations
        WHERE trip_id = $1`

	rows, err := r.db.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitations: %w", err)
	}
	defer rows.Close()

	var invitations []*Invitation
	for rows.Next() {
		var inv Invitation
		err := rows.Scan(
			&inv.ID,
			&inv.TripID,
			&inv.Email,
			&inv.Status,
			&inv.CreatedAt,
			&inv.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, &inv)
	}

	return invitations, nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM invitations WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete invitation: %w", err)
	}

	return nil
}
