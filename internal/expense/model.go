package expense

import (
	"time"
)

type Expense struct {
	ID          int64     `json:"id"`
	TripID      int64     `json:"trip_id" validate:"required"`
	Category    string    `json:"category" validate:"required,max=50"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Description string    `json:"description" validate:"max=500"`
	Date        time.Time `json:"date" validate:"required"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BudgetSummary struct {
	TotalExpenses float64            `json:"total_expenses"`
	ByCategory    map[string]float64 `json:"by_category"`
}
