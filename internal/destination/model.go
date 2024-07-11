package destination

import (
	"time"
)

type Destination struct {
	ID          int64     `json:"id"`
	TripID      int64     `json:"trip_id"`
	Name        string    `json:"name"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
