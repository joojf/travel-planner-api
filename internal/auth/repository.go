package auth

import (
	"database/sql"
	"errors"
	"time"
)

type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int64) error
	ListUsers(limit, offset int) ([]*User, error)
}

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) CreateUser(user *User) error {
	query := `
        INSERT INTO users (email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	err := r.db.QueryRow(query, user.Email, user.Password, time.Now(), time.Now()).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLRepository) GetUserByEmail(email string) (*User, error) {
	query := `
        SELECT id, email, password, created_at, updated_at
        FROM users
        WHERE email = $1`

	var user User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *SQLRepository) GetUserByID(id int64) (*User, error) {
	query := `
        SELECT id, email, password, created_at, updated_at
        FROM users
        WHERE id = $1`

	var user User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *SQLRepository) UpdateUser(user *User) error {
	query := `
        UPDATE users
        SET email = $1, password = $2, updated_at = $3
        WHERE id = $4`

	_, err := r.db.Exec(query, user.Email, user.Password, time.Now(), user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLRepository) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *SQLRepository) ListUsers(limit, offset int) ([]*User, error) {
	query := `
        SELECT id, email, created_at, updated_at
        FROM users
        ORDER BY id
        LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
