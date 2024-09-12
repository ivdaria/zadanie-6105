package entity

import (
	"github.com/google/uuid"
	"time"
)

type Feedback struct {
	ID               uuid.UUID
	BidID            uuid.UUID
	FeedbackAuthorID uuid.UUID
	Comment          string
	CreatedAt        time.Time
}
