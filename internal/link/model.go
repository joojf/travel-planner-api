package link

import (
	"time"
)

type Link struct {
	ID          int64     `json:"id"`
	TripID      int64     `json:"trip_id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
