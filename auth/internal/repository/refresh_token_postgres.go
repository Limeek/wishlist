package repository

import (
	"context"
	"database/sql"
	"errors"
	"wishlist-auth/internal/entity"

	"github.com/google/uuid"
)

type RefreshTokenRepositoryDb struct {
	db *sql.DB
}

func NewRefreshTokenRepositoryDb(db *sql.DB) *RefreshTokenRepositoryDb {
	return &RefreshTokenRepositoryDb{db: db}
}

func (r *RefreshTokenRepositoryDb) Save(ctx context.Context, token *entity.RefreshToken) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO refresh_tokens(token, user_id, expires_at) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token = $1, expires_at = $3", token.Token, token.UserId, token.ExpiresAt)
	return err
}

func (r *RefreshTokenRepositoryDb) Delete(ctx context.Context, userId uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE user_id = $1", userId)
	return err
}

func (r *RefreshTokenRepositoryDb) FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken

	row := r.db.QueryRowContext(ctx, "SELECT token, user_id, expires_at FROM refresh_tokens WHERE token = $1", token)

	err := row.Scan(
		&refreshToken.Token,
		&refreshToken.UserId,
		&refreshToken.ExpiresAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrRefreshTokenIsNotValid
	}

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}
