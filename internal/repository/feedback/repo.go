package feedback

import (
	"context"
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	db *pgx.Conn
}

func NewRepo(db *pgx.Conn) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) CreateFeedback(ctx context.Context, feedback *entity.Feedback) error {
	const query = `INSERT INTO feedback(bid_id, feedback_author, feedback) 
	VALUES($1, $2, $3)`

	mdl := modelFromFeedback(feedback)
	_, err := r.db.Exec(ctx, query, mdl.BidID, mdl.FeedbackAuthorID, mdl.Comment)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetFeedbackByTenderIDAndAuthor(ctx context.Context, tenderID uuid.UUID, authorID uuid.UUID, pagination entity.Pagination) ([]*entity.Feedback, error) {
	const query = `SELECT f.*
FROM feedback f
WHERE bid_id IN (SELECT DISTINCT ON (b.id) b.id
                 FROM bids b
                 WHERE b.tender_id = $1 AND b.creator_id = $2
                 ORDER BY b.id, b.version DESC
)
ORDER BY f.created_at DESC`

	var mdls models

	rows, err := r.db.Query(ctx, query+pagination.ToSQL(), tenderID, authorID)
	if err != nil {
		return nil, fmt.Errorf("select feedback: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mdl model
		if err := rows.Scan(&mdl); err != nil {
			return nil, fmt.Errorf("list of bids scan row:  %w", err)
		}
		mdls = append(mdls, &mdl)
	}

	return mdls.toFeedbacks(), nil
}
