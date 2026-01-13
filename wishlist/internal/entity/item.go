package entity

import "github.com/google/uuid"

type Item struct {
	Id          uuid.UUID
	Name        string
	Price       float64
	ImageUrl    string
	Link        string
	Description string
	ReservedBy  *uuid.UUID
}
