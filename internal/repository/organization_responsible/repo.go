package organization_responsible

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	db *pgx.Conn
}

func NewRepo(db *pgx.Conn) *Repo {
	return &Repo{db: db}
}

func (r *Repo) IsUserOrganizationResponsible(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (bool, error) {
	const query = `SELECT * FROM organization_responsible WHERE user_id = $1 AND organization_id = $2`

	var mdl model

	row := r.db.QueryRow(ctx, query, userID, orgID)
	if err := row.Scan(&mdl); err != nil {
		return false, err
	}

	return true, nil
}
