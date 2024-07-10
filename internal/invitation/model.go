package invitation

import (
	"time"
)

type Invitation struct {
	ID        int64     `json:"id"`
	TripID    int64     `json:"trip_id"`
	Email     string    `json:"email"`
	Status    string    `json:"status"` // "pending", "accepted", "rejected"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
