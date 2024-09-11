package tenders

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

func (r *Repo) CreateTender(ctx context.Context, tender *entity.Tender) (*entity.Tender, error) {
	const query = `INSERT INTO tenders(name, description, service_type, status, organization_id, creator_user_id, version) 
	VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING *`
	mdl := modelFromTender(tender)
	err := r.db.QueryRow(ctx, query, mdl.Name, mdl.Description, mdl.ServiceType, mdl.Status, mdl.OrganizationID, mdl.CreatorID, mdl.Version).
		Scan(&mdl.ID, &mdl.Name, &mdl.Description, &mdl.ServiceType, &mdl.Status, &mdl.OrganizationID, &mdl.CreatorID, &mdl.Version, &mdl.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create tender: %w", err)
	}

	newTender := mdl.toTender()
	return newTender, nil
}

func (r *Repo) GetAllTenders(ctx context.Context, filter entity.GetTendersFilter, pagination entity.Pagination) ([]*entity.Tender, error) {
	params := pgx.NamedArgs{}
	query := `SELECT DISTINCT ON (id) * FROM tenders`

	if filter.ServiceTypes != nil {
		query += ` WHERE service_type = ANY(@serviceTypes)`
		params["serviceTypes"] = *filter.ServiceTypes
	}

	query += ` ORDER BY id, tenders.version DESC, name ASC`
	query += pagination.ToSQL()

	var mdls models

	rows, err := r.db.Query(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("select all tenders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mdl model
		if err := rows.Scan(&mdl); err != nil {
			return nil, fmt.Errorf("list of tenders scan row:  %w", err)
		}
		mdls = append(mdls, &mdl)
	}

	return mdls.toTenders(), nil
}

func (r *Repo) GetTendersByUsername(ctx context.Context, username string, pagination entity.Pagination) ([]*entity.Tender, error) {
	const query = `
SELECT DISTINCT ON (t.id) t.*
FROM tenders t
INNER JOIN employee e ON t.creator_user_id = e.id
WHERE e.username = @username
ORDER BY t.id, t.version DESC, name
`

	var mdls models

	rows, err := r.db.Query(ctx, query+pagination.ToSQL(), pgx.NamedArgs{"username": username})
	if err != nil {
		return nil, fmt.Errorf("select all tenders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mdl model
		if err := rows.Scan(&mdl); err != nil {
			return nil, fmt.Errorf("list of tenders scan row:  %w", err)
		}
		mdls = append(mdls, &mdl)
	}

	return mdls.toTenders(), nil
}

func (r *Repo) GetTenderByID(ctx context.Context, id uuid.UUID) (*entity.Tender, error) {
	const query = `SELECT * FROM tenders WHERE id = $1 ORDER BY version DESC LIMIT 1`

	var mdl model

	row := r.db.QueryRow(ctx, query, id)
	if err := row.Scan(&mdl); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return mdl.toTender(), nil
}

func (r *Repo) GetTenderByIDAndVersion(ctx context.Context, id uuid.UUID, version int) (*entity.Tender, error) {
	const query = `SELECT * FROM tenders WHERE id = $1 AND version = $2`

	var mdl model

	row := r.db.QueryRow(ctx, query, id, version)
	if err := row.Scan(&mdl); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return mdl.toTender(), nil
}

func (r *Repo) UpdateTender(ctx context.Context, tender *entity.Tender) (*entity.Tender, error) {
	const query = `INSERT INTO tenders(id, name, description, service_type, status, organization_id, creator_user_id, version, created_at) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *`

	mdl := modelFromTender(tender)
	err := r.db.QueryRow(ctx, query, mdl.ID, mdl.Name, mdl.Description, mdl.ServiceType, mdl.Status, mdl.OrganizationID, mdl.CreatorID, mdl.Version, mdl.CreatedAt).
		Scan(&mdl.ID, &mdl.Name, &mdl.Description, &mdl.ServiceType, &mdl.Status, &mdl.OrganizationID, &mdl.CreatorID, &mdl.Version, &mdl.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("update tender: %w", err)
	}

	return mdl.toTender(), nil
}

func (r *Repo) UpdateTenderStatus(ctx context.Context, id uuid.UUID, newStatus entity.TenderStatus) error {
	const query = `UPDATE tenders SET status = $2 WHERE id = $1`

	commandTag, err := r.db.Exec(ctx, query, id, newStatus)
	if err != nil {
		return fmt.Errorf("update tender status: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return er.ErrNoRowsAffected
	}

	return nil
}
