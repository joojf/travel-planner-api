package review

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

func (r *Repository) Create(review *Review) error {
	query := `
        INSERT INTO reviews (trip_id, user_id, activity_id, place_id, rating, comment, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		review.TripID,
		review.UserID,
		review.ActivityID,
		review.Rating,
		review.Comment,
		time.Now(),
		time.Now(),
	).Scan(&review.ID)

	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

func (r *Repository) GetByTripID(tripID int64) ([]*Review, error) {
	query := `
        SELECT id, trip_id, user_id, activity_id, place_id, rating, comment, created_at, updated_at
        FROM reviews
        WHERE trip_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	defer rows.Close()

	var reviews []*Review
	for rows.Next() {
		var review Review
		err := rows.Scan(
			&review.ID,
			&review.TripID,
			&review.UserID,
			&review.ActivityID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		reviews = append(reviews, &review)
	}

	return reviews, nil
}

func (r *Repository) GetByID(id int64) (*Review, error) {
	query := `
        SELECT id, trip_id, user_id, activity_id, place_id, rating, comment, created_at, updated_at
        FROM reviews
        WHERE id = $1`

	var review Review
	err := r.db.QueryRow(query, id).Scan(
		&review.ID,
		&review.TripID,
		&review.UserID,
		&review.ActivityID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return &review, nil
}

func (r *Repository) Update(review *Review) error {
	query := `
        UPDATE reviews
        SET rating = $1, comment = $2, updated_at = $3
        WHERE id = $4 AND user_id = $5`

	_, err := r.db.Exec(
		query,
		review.Rating,
		review.Comment,
		time.Now(),
		review.ID,
		review.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}

func (r *Repository) Delete(id, userID int64) error {
	query := `DELETE FROM reviews WHERE id = $1 AND user_id = $2`

	_, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}
