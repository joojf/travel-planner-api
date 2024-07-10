package expense

import (
	"time"
)

type Expense struct {
	ID          int64     `json:"id"`
	TripID      int64     `json:"trip_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	PaidBy      int64     `json:"paid_by"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
