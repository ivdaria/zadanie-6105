package tenders

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type model struct {
	ID             uuid.UUID
	Name           string
	Description    string
	ServiceType    entity.ServiceType
	Status         entity.TenderStatus
	OrganizationID uuid.UUID
	CreatorID      uuid.UUID
	Version        int
	CreatedAt      time.Time
}

type models []*model

func (m *model) ScanRow(rows pgx.Rows) error {
	return rows.Scan(&m.ID, &m.Name, &m.Description, &m.ServiceType, &m.Status, &m.OrganizationID, &m.CreatorID, &m.Version, &m.CreatedAt)
}

func (m *model) toTender() *entity.Tender {
	return &entity.Tender{
		ID:             m.ID,
		Name:           m.Name,
		Description:    m.Description,
		ServiceType:    m.ServiceType,
		Status:         m.Status,
		OrganizationID: m.OrganizationID,
		CreatorID:      m.CreatorID,
		Version:        m.Version,
		CreatedAt:      m.CreatedAt,
	}
}

func (mdls models) toTenders() []*entity.Tender {
	if len(mdls) == 0 {
		return nil
	}

	result := make([]*entity.Tender, 0, len(mdls))
	for _, m := range mdls {
		result = append(result, m.toTender())
	}

	return result
}

func modelFromTender(item *entity.Tender) *model {
	return &model{
		ID:             item.ID,
		Name:           item.Name,
		Description:    item.Description,
		ServiceType:    item.ServiceType,
		Status:         item.Status,
		OrganizationID: item.OrganizationID,
		CreatorID:      item.CreatorID,
		Version:        item.Version,
		CreatedAt:      item.CreatedAt,
	}
}
