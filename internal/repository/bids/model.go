package bids

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type model struct {
	ID          uuid.UUID
	Name        string
	TenderID    uuid.UUID
	CreatorID   uuid.UUID
	Description string
	Decision    entity.BidDecision
	Status      entity.BidStatus
	AuthorType  entity.BidAuthorType
	Version     int
	CreatedAt   time.Time
}

type models []*model

func (m *model) ScanRow(rows pgx.Rows) error {
	return rows.Scan(&m.ID, &m.TenderID, &m.CreatorID, &m.Name, &m.Description, &m.Decision, &m.Status, &m.AuthorType, &m.Version, &m.CreatedAt)
}

func (m *model) toBid() *entity.Bid {
	return &entity.Bid{
		ID:          m.ID,
		Name:        m.Name,
		TenderID:    m.TenderID,
		CreatorID:   m.CreatorID,
		Description: m.Description,
		Decision:    m.Decision,
		Status:      m.Status,
		AuthorType:  m.AuthorType,
		Version:     m.Version,
		CreatedAt:   m.CreatedAt,
	}
}

func (mdls models) toBids() []*entity.Bid {
	if len(mdls) == 0 {
		return nil
	}

	result := make([]*entity.Bid, 0, len(mdls))
	for _, m := range mdls {
		result = append(result, m.toBid())
	}

	return result
}

func modelFromBid(item *entity.Bid) *model {
	return &model{
		ID:          item.ID,
		Name:        item.Name,
		TenderID:    item.TenderID,
		CreatorID:   item.CreatorID,
		Description: item.Description,
		Decision:    item.Decision,
		Status:      item.Status,
		AuthorType:  item.AuthorType,
		Version:     item.Version,
		CreatedAt:   item.CreatedAt,
	}
}
