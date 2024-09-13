package organization

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type model struct {
	ID          uuid.UUID
	Name        string
	Description string
	Type        entity.OrganizationType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//nolint:unused
type models []*model

func (m *model) ScanRow(rows pgx.Rows) error {
	return rows.Scan(&m.ID, &m.Name, &m.Description, &m.Type, &m.CreatedAt, &m.UpdatedAt)
}

func (m *model) toOrganization() *entity.Organization {
	return &entity.Organization{
		ID:        m.ID,
		Name:      m.Name,
		Type:      m.Type,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

//nolint:unused
func (mdls models) toOrganizations() []*entity.Organization {
	if len(mdls) == 0 {
		return nil
	}

	result := make([]*entity.Organization, 0, len(mdls))
	for _, m := range mdls {
		result = append(result, m.toOrganization())
	}

	return result
}

//nolint:unused
func modelFromOrganization(item *entity.Organization) *model {
	return &model{
		ID:        item.ID,
		Name:      item.Name,
		Type:      item.Type,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
