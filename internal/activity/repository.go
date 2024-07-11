package activity

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
	Create(activity *Activity) error
	GetByTripID(tripID int64) ([]*Activity, error)
	GetByID(id int64) (*Activity, error)
	Update(activity *Activity) error
	Delete(id int64) error
}

var _ RepositoryInterface = (*Repository)(nil)

func (r *Repository) Create(activity *Activity) error {
	query := `
        INSERT INTO activities (trip_id, name, description, location, start_time, end_time, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		activity.TripID,
		activity.Name,
		activity.Description,
		activity.Location,
		activity.StartTime,
		activity.EndTime,
		time.Now(),
		time.Now(),
	).Scan(&activity.ID)

	if err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}

	return nil
}

func (r *Repository) GetByTripID(tripID int64) ([]*Activity, error) {
	query := `
        SELECT id, trip_id, name, description, location, start_time, end_time, created_at, updated_at
        FROM activities
        WHERE trip_id = $1`

	rows, err := r.db.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}
	defer rows.Close()

	var activities []*Activity
	for rows.Next() {
		var a Activity
		err := rows.Scan(
			&a.ID,
			&a.TripID,
			&a.Name,
			&a.Description,
			&a.Location,
			&a.StartTime,
			&a.EndTime,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, &a)
	}

	return activities, nil
}

func (r *Repository) GetByID(id int64) (*Activity, error) {
	query := `
        SELECT id, trip_id, name, description, location, start_time, end_time, created_at, updated_at
        FROM activities
        WHERE id = $1`

	var activity Activity
	err := r.db.QueryRow(query, id).Scan(
		&activity.ID,
		&activity.TripID,
		&activity.Name,
		&activity.Description,
		&activity.Location,
		&activity.StartTime,
		&activity.EndTime,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("activity not found")
		}
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	return &activity, nil
}

func (r *Repository) Update(activity *Activity) error {
	query := `
        UPDATE activities
        SET name = $1, description = $2, location = $3, start_time = $4, end_time = $5, updated_at = $6
        WHERE id = $7`

	_, err := r.db.Exec(
		query,
		activity.Name,
		activity.Description,
		activity.Location,
		activity.StartTime,
		activity.EndTime,
		time.Now(),
		activity.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update activity: %w", err)
	}

	return nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM activities WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete activity: %w", err)
	}

	return nil
}
