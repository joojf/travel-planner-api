package itinerary

import (
	"time"
)

type Itinerary struct {
	ID          int64     `json:"id"`
	TripID      int64     `json:"trip_id" validate:"required"`
	Title       string    `json:"title" validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"max=500"`
	Date        time.Time `json:"date" validate:"required"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
