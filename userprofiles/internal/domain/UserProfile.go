package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserId    uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Avatar    string
	BirthDate time.Time
}
