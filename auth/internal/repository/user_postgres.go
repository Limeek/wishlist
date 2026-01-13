package repository

import (
	"context"
	"database/sql"
	"errors"
	"wishlist-auth/internal/entity"
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) UserRepository {
	return &UserRepositoryPostgres{
		db: db,
	}
}

func (r *UserRepositoryPostgres) Save(ctx context.Context, user *entity.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO users (id, email, password) VALUES ($1, $2, $3)",
		user.ID, user.Email, user.Password)

	return err
}

func (r *UserRepositoryPostgres) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
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
		return nil, entity.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
