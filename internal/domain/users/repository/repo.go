package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"icu/internal/domain/users/models"
)

type UserRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUserRepository(db *sql.DB, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) UserExists(ctx context.Context, email string) (int, error) {
	query := `SELECT 1 from users where email = $1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, nil
	case err != nil:
		return 0, fmt.Errorf("error checking if user exists: %s", err)
	default:
		return 1, nil
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	// Обратите внимание, что `created_at` и `updated_at` убраны из запроса, так как они устанавливаются автоматически
	query := `INSERT INTO users(email, username, password, role) VALUES($1, $2, $3, $4) RETURNING email`
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Username, user.Password, user.Role).Scan(&user.Email)
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}
	r.logger.Info("data inserted in table 'users'")

	return nil
}
func (r *UserRepository) FindAll(ctx context.Context) ([]*models.User, error) {
	query := `SELECT email, username, password, role, created_at, updated_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve users: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Email, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan row: %v", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return users, nil
}

func (r *UserRepository) FindOne(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT email, username, password, role, created_at, updated_at FROM users WHERE email = $1`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.Email, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("could not retrieve user: %v", err)
	}
	return &user, nil
}

func (r *UserRepository) DeleteByEmail(ctx context.Context, email string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not start transaction: %v", err)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM users WHERE email = $1`, email)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not delete user: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	r.logger.Info(fmt.Sprintf("user %s was removed from table 'users'", email))
	return nil
}
