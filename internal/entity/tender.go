package entity

import (
	"github.com/google/uuid"
	"time"
)

type ServiceType string

const (
	ServiceTypeConstruction ServiceType = "Construction"
	ServiceTypeDelivery     ServiceType = "Delivery"
	ServiceTypeManufacture  ServiceType = "Manufacture"
)

type TenderStatus string

const (
	TenderStatusCreated   TenderStatus = "Created"
	TenderStatusPublished TenderStatus = "Published"
	TenderStatusClosed    TenderStatus = "Closed"
)

type Tender struct {
	ID             uuid.UUID
	Name           string
	Description    string
	ServiceType    ServiceType
	Status         TenderStatus
	OrganizationID uuid.UUID
	CreatorID      uuid.UUID
	Version        int
	CreatedAt      time.Time
}

type GetTendersFilter struct {
	ServiceTypes *[]string
}
