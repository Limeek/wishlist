package repository

import (
	"context"
	"database/sql"
	"errors"
	"userprofiles/internal/domain"

	"github.com/google/uuid"
)

type UserProfileRepositoryDb struct {
	db *sql.DB
}

func NewUserProfileRepositoryDb(db *sql.DB) *UserProfileRepositoryDb {
	return &UserProfileRepositoryDb{
		db: db,
	}
}

func (r *UserProfileRepositoryDb) Create(ctx context.Context, userProfile domain.UserProfile) error {
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

func (r *UserProfileRepositoryDb) Get(ctx context.Context, userId uuid.UUID) (*domain.UserProfile, error) {
	var userProfile domain.UserProfile

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
		return nil, domain.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &userProfile, nil
}
