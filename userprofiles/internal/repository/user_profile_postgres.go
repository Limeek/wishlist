package repository

import (
	"context"
	"database/sql"
	"errors"
	"userprofiles/internal/entity"

	"github.com/google/uuid"
)

type UserProfileRepositoryPostgres struct {
	db *sql.DB
}

func NewUserProfileRepositoryPostgres(db *sql.DB) *UserProfileRepositoryPostgres {
	return &UserProfileRepositoryPostgres{
		db: db,
	}
}

func (r *UserProfileRepositoryPostgres) Create(ctx context.Context, userProfile entity.UserProfile) error {
	_, err := r.db.ExecContext(ctx, "insert into user_profiles(user_id, first_name, last_name, email, phone, avatar, birth_date) values ($1, $2 ,$3 ,$4, $5, $6, $7) ON CONFLICT (user_id) DO NOTHING;",
		userProfile.UserId,
		userProfile.FirstName,
		userProfile.LastName,
		userProfile.Email,
		userProfile.Phone,
		userProfile.Avatar,
		userProfile.BirthDate)

	return err
}

func (r *UserProfileRepositoryPostgres) Get(ctx context.Context, userId uuid.UUID) (*entity.UserProfile, error) {
	var userProfile entity.UserProfile

	row := r.db.QueryRowContext(ctx, "select user_id, first_name, last_name, email, phone, avatar, birth_date from user_profiles where user_id = $1", userId)

	err := row.Scan(
		&userProfile.UserId,
		&userProfile.FirstName,
		&userProfile.LastName,
		&userProfile.Email,
		&userProfile.Phone,
		&userProfile.Avatar,
		&userProfile.BirthDate,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, entity.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &userProfile, nil
}
