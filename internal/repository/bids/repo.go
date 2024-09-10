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
	const query = `INSERT INTO bids(tender_id, creator_id, organization_id, decision, status, author_type, version, created_at) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	var id uuid.UUID
	mdl := modelFromBid(bid)
	err := r.db.QueryRow(ctx, query, mdl.TenderID, mdl.CreatorID, mdl.OrganizationID, mdl.Decision, mdl.Status, mdl.AuthorType, mdl.Version, mdl.CreatedAt).Scan(&id)
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

func (r *Repo) GetBidsByCreatorID(ctx context.Context, id uuid.UUID) ([]*entity.Bid, error) {
	const query = `SELECT * FROM bids WHERE creator_id = $1`

	var mdls models

	rows, err := r.db.Query(ctx, query, id)
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

func (r *Repo) GetBidsByTenderID(ctx context.Context, id uuid.UUID) ([]*entity.Bid, error) {
	const query = `SELECT * FROM bids WHERE tender_id = $1`

	var mdls models

	rows, err := r.db.Query(ctx, query, id)
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

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Bid, error) {
	const query = `SELECT * FROM bids WHERE tender_id = $1`

	var mdl model

	row := r.db.QueryRow(ctx, query, id)
	if err := row.Scan(&mdl); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return mdl.toBid(), nil
}

func (r *Repo) UpdateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error) {
	const query = `UPDATE bids SET tender_id = $2, creator_id = $3, organization_id = $4, decision = $5, status = $6, author_type = $7, version = $8 WHERE id = $1 RETURNING *`

	mdl := modelFromBid(bid)

	commandTag, err := r.db.Exec(ctx, query, mdl.ID, mdl.TenderID, mdl.CreatorID, mdl.OrganizationID, mdl.Decision, mdl.Status, mdl.AuthorType, mdl.Version)

	if err != nil {
		return nil, fmt.Errorf("update bids: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return nil, er.ErrNoRowsAffected
	}

	return mdl.toBid(), nil
}
