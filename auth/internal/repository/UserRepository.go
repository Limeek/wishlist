package repository

import (
	"context"
	"database/sql"
	"errors"
	"wishlist-auth/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (id, email, password) VALUES ($1, $2, $3)",
		user.ID, user.Email, user.Password)

	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	row := r.db.QueryRowContext(ctx, `
        SELECT id, email, password
        FROM users WHERE email = $1
    `, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
