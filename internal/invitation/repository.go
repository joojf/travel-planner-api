package invitation

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(inv *Invitation) error {
	query := `INSERT INTO invitations (trip_id, email, status, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5)
              RETURNING id`

	err := r.db.QueryRow(query, inv.TripID, inv.Email, "pending", time.Now(), time.Now()).Scan(&inv.ID)
	return err
}

func (r *Repository) GetByTripID(tripID int64) ([]Invitation, error) {
	query := `SELECT id, trip_id, email, status, created_at, updated_at
              FROM invitations
              WHERE trip_id = $1`

	rows, err := r.db.Query(query, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []Invitation
	for rows.Next() {
		var inv Invitation
		err := rows.Scan(&inv.ID, &inv.TripID, &inv.Email, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, inv)
	}
	return invitations, nil
}

// TODO: Implement other necessary methods (Update, Delete, etc.)
