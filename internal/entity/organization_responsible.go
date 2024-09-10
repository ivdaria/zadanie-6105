package entity

import "github.com/google/uuid"

type OrganizationResponsible struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}
