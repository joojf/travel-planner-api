package review

import (
	"time"
)

type Review struct {
	ID         int64     `json:"id"`
	TripID     int64     `json:"trip_id" validate:"required"`
	UserID     int64     `json:"user_id" validate:"required"`
	ActivityID *int64    `json:"activity_id,omitempty"`
	Rating     int       `json:"rating" validate:"required,min=1,max=5"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
