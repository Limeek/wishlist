package entity

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
}

type UserDetails struct {
	Phone     string `json:"phoneNumber"`
	BirthDate string `json:"birthDate"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarURL string `json:"avatarUrl"`
}
