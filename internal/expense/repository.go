package expense

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
	Create(expense *Expense) error
	GetByTripID(tripID int64) ([]*Expense, error)
	GetByID(id int64) (*Expense, error)
	Update(expense *Expense) error
	Delete(id int64) error
	GetBudgetSummary(tripID int64) (*BudgetSummary, error)
}

func (r *Repository) Create(expense *Expense) error {
	query := `
        INSERT INTO expenses (trip_id, category, amount, description, date, created_by, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`

	err := r.db.QueryRow(
		query,
		expense.TripID,
		expense.Category,
		expense.Amount,
		expense.Description,
		expense.Date,
		expense.CreatedBy,
		time.Now(),
		time.Now(),
	).Scan(&expense.ID)

	if err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}

	return nil
}

func (r *Repository) GetByTripID(tripID int64) ([]*Expense, error) {
	query := `
        SELECT id, trip_id, category, amount, description, date, created_by, created_at, updated_at
        FROM expenses
        WHERE trip_id = $1
        ORDER BY date DESC`

	rows, err := r.db.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses: %w", err)
	}
	defer rows.Close()

	var expenses []*Expense
	for rows.Next() {
		var expense Expense
		err := rows.Scan(
			&expense.ID,
			&expense.TripID,
			&expense.Category,
			&expense.Amount,
			&expense.Description,
			&expense.Date,
			&expense.CreatedBy,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expense: %w", err)
		}
		expenses = append(expenses, &expense)
	}

	return expenses, nil
}

func (r *Repository) GetByID(id int64) (*Expense, error) {
	query := `
        SELECT id, trip_id, category, amount, description, date, created_by, created_at, updated_at
        FROM expenses
        WHERE id = $1`

	var expense Expense
	err := r.db.QueryRow(query, id).Scan(
		&expense.ID,
		&expense.TripID,
		&expense.Category,
		&expense.Amount,
		&expense.Description,
		&expense.Date,
		&expense.CreatedBy,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("expense not found")
		}
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}

	return &expense, nil
}

func (r *Repository) Update(expense *Expense) error {
	query := `
        UPDATE expenses
        SET category = $1, amount = $2, description = $3, date = $4, updated_at = $5
        WHERE id = $6`

	_, err := r.db.Exec(
		query,
		expense.Category,
		expense.Amount,
		expense.Description,
		expense.Date,
		time.Now(),
		expense.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	return nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM expenses WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	return nil
}

func (r *Repository) GetBudgetSummary(tripID int64) (*BudgetSummary, error) {
	query := `
        SELECT category, SUM(amount) as total
        FROM expenses
        WHERE trip_id = $1
        GROUP BY category`

	rows, err := r.db.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget summary: %w", err)
	}
	defer rows.Close()

	summary := &BudgetSummary{
		ByCategory: make(map[string]float64),
	}

	for rows.Next() {
		var category string
		var amount float64
		err := rows.Scan(&category, &amount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget summary: %w", err)
		}
		summary.ByCategory[category] = amount
		summary.TotalExpenses += amount
	}

	return summary, nil
}
