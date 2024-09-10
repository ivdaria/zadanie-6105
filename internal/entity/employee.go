package entity

import (
	"github.com/google/uuid"
	"time"
)

type Employee struct {
	ID        uuid.UUID
	UserName  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
