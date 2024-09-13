package organization_responsible

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type model struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}

//nolint:unused
type models []*model

func (m *model) ScanRow(rows pgx.Rows) error {
	return rows.Scan(&m.ID, &m.OrganizationID, &m.EmployeeID)
}

func (m *model) toOrganizationResponsible() *entity.OrganizationResponsible {
	return &entity.OrganizationResponsible{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		EmployeeID:     m.EmployeeID,
	}
}

//nolint:unused
func (mdls models) toOrganizations() []*entity.OrganizationResponsible {
	if len(mdls) == 0 {
		return nil
	}

	result := make([]*entity.OrganizationResponsible, 0, len(mdls))
	for _, m := range mdls {
		result = append(result, m.toOrganizationResponsible())
	}

	return result
}

//nolint:unused
func modelFromOrganizationResponsible(item *entity.OrganizationResponsible) *model {
	return &model{
		ID:             item.ID,
		OrganizationID: item.OrganizationID,
		EmployeeID:     item.EmployeeID,
	}
}
