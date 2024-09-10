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
	BidDecisionApproved BidDecision = "Approved"
	BidDecisionRejected BidDecision = "Rejected"
)

type BidStatus string

const (
	BidStatusApproved  BidStatus = "Approved"
	BidStatusCanceled  BidStatus = "Canceled"
	BidStatusCreated   BidStatus = "Created"
	BidStatusPublished BidStatus = "Published"
	BidStatusRejected  BidStatus = "Rejected"
)

type Bid struct {
	ID             uuid.UUID
	TenderID       uuid.UUID
	CreatorID      uuid.UUID
	OrganizationID uuid.UUID
	Decision       BidDecision
	Status         BidStatus
	AuthorType     BidAuthorType
	Version        int
	CreatedAt      time.Time
}
