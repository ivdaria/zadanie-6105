package entity

import (
	"github.com/google/uuid"
)

type Feedback struct {
	ID               uuid.UUID
	BidID            uuid.UUID
	BidAuthorID      uuid.UUID
	TenderID         uuid.UUID
	FeedbackAuthorID uuid.UUID
	Comment          string
}
