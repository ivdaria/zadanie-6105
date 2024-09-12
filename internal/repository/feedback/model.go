package feedback

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type model struct {
	ID               uuid.UUID
	BidID            uuid.UUID
	BidAuthorID      uuid.UUID
	TenderID         uuid.UUID
	FeedbackAuthorID uuid.UUID
	Comment          string
}

type models []*model

func (m *model) ScanRow(rows pgx.Rows) error {
	return rows.Scan(&m.ID, &m.BidID, &m.BidAuthorID, &m.TenderID, &m.FeedbackAuthorID, &m.Comment)
}

func (m *model) toFeedback() *entity.Feedback {
	return &entity.Feedback{
		ID:               m.ID,
		BidID:            m.BidID,
		BidAuthorID:      m.BidAuthorID,
		TenderID:         m.TenderID,
		FeedbackAuthorID: m.FeedbackAuthorID,
		Comment:          m.Comment,
	}
}

func (mdls models) toFeedbacks() []*entity.Feedback {
	if len(mdls) == 0 {
		return nil
	}

	result := make([]*entity.Feedback, 0, len(mdls))
	for _, m := range mdls {
		result = append(result, m.toFeedback())
	}

	return result
}

func modelFromFeedback(item *entity.Feedback) *model {
	return &model{
		ID:               item.ID,
		BidID:            item.BidID,
		BidAuthorID:      item.BidAuthorID,
		TenderID:         item.TenderID,
		FeedbackAuthorID: item.FeedbackAuthorID,
		Comment:          item.Comment,
	}
}
