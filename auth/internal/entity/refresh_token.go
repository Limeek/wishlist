package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string
	UserId    uuid.UUID
	ExpiresAt time.Time
}
