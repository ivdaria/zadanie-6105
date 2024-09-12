package bids

import (
	"context"
	"fmt"
	er "git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/errors"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	db *pgx.Conn
}

func NewRepo(db *pgx.Conn) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error) {
	const query = `INSERT INTO bids(name, tender_id, creator_id, description, decision, status, author_type, version) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`

	mdl := modelFromBid(bid)
	err := r.db.QueryRow(ctx, query, mdl.Name, mdl.TenderID, mdl.CreatorID, mdl.Description, mdl.Decision, mdl.Status, mdl.AuthorType, mdl.Version).Scan(
		&mdl.ID,
		&mdl.TenderID,
		&mdl.CreatorID,
		&mdl.Name,
		&mdl.Description,
		&mdl.Decision,
		&mdl.Status,
		&mdl.AuthorType,
		&mdl.Version,
		&mdl.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create bid: %w", err)
	}

	newBid := mdl.toBid()
	return newBid, nil
}

func (r *Repo) GetAllBids(ctx context.Context) ([]*entity.Bid, error) {
	const query = `SELECT * FROM bids`

	var mdls models

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select all bids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mdl model
		if err := rows.Scan(&mdl); err != nil {
			return nil, fmt.Errorf("list of bids scan row:  %w", err)
		}
		mdls = append(mdls, &mdl)
	}

	return mdls.toBids(), nil
}

func (r *Repo) GetBidsByUsername(ctx context.Context, username string, pagination entity.Pagination) ([]*entity.Bid, error) {
	const query = `
SELECT DISTINCT ON (b.id) b.*
FROM bids b
INNER JOIN employee e ON b.creator_id = e.id
WHERE e.username = @username
ORDER BY b.id, b.version DESC, name
`

	var mdls models

	rows, err := r.db.Query(ctx, query+pagination.ToSQL(), pgx.NamedArgs{"username": username})
	if err != nil {
		return nil, fmt.Errorf("select all bids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mdl model
		if err := rows.Scan(&mdl); err != nil {
			return nil, fmt.Errorf("list of bids scan row:  %w", err)
		}
		mdls = append(mdls, &mdl)
	}

	return mdls.toBids(), nil
}

func (r *Repo) GetBidsByTenderID(ctx context.Context, tenderID uuid.UUID, pagination entity.Pagination) ([]*entity.Bid, error) {
	const query = `
SELECT DISTINCT ON (b.id) b.*
FROM bids b
WHERE b.tender_id = @tender_id AND status = ANY('{Published,Canceled}')
ORDER BY b.id, b.version DESC, name
`

	var mdls models

	rows, err := r.db.Query(ctx, query+pagination.ToSQL(), pgx.NamedArgs{"tender_id": tenderID})
	if err != nil {
		return nil, fmt.Errorf("select all bids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mdl model
		if err := rows.Scan(&mdl); err != nil {
			return nil, fmt.Errorf("list of bids scan row:  %w", err)
		}
		mdls = append(mdls, &mdl)
	}

	return mdls.toBids(), nil
}

func (r *Repo) GetBidByID(ctx context.Context, id uuid.UUID) (*entity.Bid, error) {
	const query = `SELECT * FROM bids WHERE id = $1 ORDER BY version DESC LIMIT 1`

	var mdl model

	row := r.db.QueryRow(ctx, query, id)
	if err := row.Scan(&mdl); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return mdl.toBid(), nil
}

func (r *Repo) GetBidByIDAndVersion(ctx context.Context, id uuid.UUID, version int) (*entity.Bid, error) {
	const query = `SELECT * FROM bids WHERE id = $1 AND version = $2`

	var mdl model

	row := r.db.QueryRow(ctx, query, id, version)
	if err := row.Scan(&mdl); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return mdl.toBid(), nil
}

func (r *Repo) UpdateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error) {
	const query = `INSERT INTO bids(id, tender_id, creator_id, name, description, decision, status, author_type, version, created_at) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *`

	mdl := modelFromBid(bid)
	err := r.db.QueryRow(ctx, query, mdl.ID, mdl.TenderID, mdl.CreatorID, mdl.Name, mdl.Description, mdl.Decision, mdl.Status, mdl.AuthorType, mdl.Version, mdl.CreatedAt).
		Scan(&mdl.ID, &mdl.TenderID, &mdl.CreatorID, &mdl.Name, &mdl.Description, &mdl.Decision, &mdl.Status, &mdl.AuthorType, &mdl.Version, &mdl.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("update bid: %w", err)
	}

	return mdl.toBid(), nil
}

func (r *Repo) UpdateBidDecision(ctx context.Context, id uuid.UUID, bidDecision entity.BidDecision) error {
	const query = `UPDATE bids SET decision = $2 WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id, bidDecision)
	if err != nil {
		return fmt.Errorf("update bid decision: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return er.ErrNoRowsAffected
	}

	return nil
}

func (r *Repo) UpdateBidStatus(ctx context.Context, id uuid.UUID, newStatus entity.BidStatus) error {
	const query = `UPDATE bids SET status = $2 WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id, newStatus)
	if err != nil {
		return fmt.Errorf("update bid status: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return er.ErrNoRowsAffected
	}

	return nil
}
