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

func (t *Tender) Patch(newName *string, newDescription *string, newServiceType *ServiceType) *Tender {
	patchedTender := *t
	if newName != nil {
		patchedTender.Name = *newName
	}
	if newDescription != nil {
		patchedTender.Description = *newDescription
	}
	if newServiceType != nil {
		patchedTender.ServiceType = *newServiceType
	}

	patchedTender.Version += 1

	return &patchedTender
}

func (t *Tender) Rollback(tenderToRollback *Tender) *Tender {
	result := *tenderToRollback
	result.Version = t.Version + 1

	return &result
}

type GetTendersFilter struct {
	ServiceTypes *[]string
}
