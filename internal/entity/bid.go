package entity

import (
	"github.com/google/uuid"
	"time"
)

type BidAuthorType string

const (
	BidAuthorTypeOrganization BidAuthorType = "Organization"
	BidAuthorTypeUser         BidAuthorType = "User"
)

type BidDecision string

const (
	BidDecisionNone     BidDecision = "No decision"
	BidDecisionApproved BidDecision = "Approved"
	BidDecisionRejected BidDecision = "Rejected"
)

type BidStatus string

const (
	BidStatusCanceled  BidStatus = "Canceled"
	BidStatusCreated   BidStatus = "Created"
	BidStatusPublished BidStatus = "Published"
)

type Bid struct {
	ID          uuid.UUID
	Name        string
	TenderID    uuid.UUID
	CreatorID   uuid.UUID
	Description string
	Decision    BidDecision
	Status      BidStatus
	AuthorType  BidAuthorType
	Version     int
	CreatedAt   time.Time
}

func (b *Bid) Patch(newName *string, newDescription *string) *Bid {
	patchedBid := *b
	if newName != nil {
		patchedBid.Name = *newName
	}
	if newDescription != nil {
		patchedBid.Description = *newDescription
	}
	patchedBid.Version += 1

	return &patchedBid
}

func (b *Bid) Rollback(bidToRollback *Bid) *Bid {
	result := *bidToRollback
	result.Version = b.Version + 1

	return &result
}
