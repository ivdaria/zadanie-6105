package feedback

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type model struct {
	ID               uuid.UUID
	BidID            uuid.UUID
	FeedbackAuthorID uuid.UUID
	Comment          string
	CreatedAt        time.Time
}

type models []*model

func (m *model) ScanRow(rows pgx.Rows) error {
	return rows.Scan(&m.ID, &m.BidID, &m.FeedbackAuthorID, &m.Comment, &m.CreatedAt)
}

func (m *model) toFeedback() *entity.Feedback {
	return &entity.Feedback{
		ID:               m.ID,
		BidID:            m.BidID,
		FeedbackAuthorID: m.FeedbackAuthorID,
		Comment:          m.Comment,
		CreatedAt:        m.CreatedAt,
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
		FeedbackAuthorID: item.FeedbackAuthorID,
		Comment:          item.Comment,
		CreatedAt:        item.CreatedAt,
	}
}
