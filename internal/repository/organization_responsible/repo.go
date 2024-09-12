package organization_responsible

import (
	"context"
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

func (r *Repo) IsUserOrganizationResponsible(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) error {
	const query = `SELECT * FROM organization_responsible WHERE user_id = $1 AND organization_id = $2`

	var mdl model

	row := r.db.QueryRow(ctx, query, userID, orgID)
	if err := row.Scan(&mdl); err != nil {
		return err
	}

	return nil
}

func (r *Repo) IsUserResponsible(ctx context.Context, userID uuid.UUID) error {
	const query = `SELECT * FROM organization_responsible WHERE user_id = $1`

	var mdl model

	row := r.db.QueryRow(ctx, query, userID)
	if err := row.Scan(&mdl); err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetOrganizationResponsibleByUserID(ctx context.Context, userID uuid.UUID) (*entity.OrganizationResponsible, error) {
	const query = `SELECT * FROM organization_responsible WHERE user_id = $1`

	var mdl model

	row := r.db.QueryRow(ctx, query, userID)
	if err := row.Scan(&mdl.ID, &mdl.OrganizationID, &mdl.EmployeeID); err != nil {
		return nil, err
	}
	return mdl.toOrganizationResponsible(), nil
}
