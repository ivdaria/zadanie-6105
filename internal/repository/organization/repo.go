package organization

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
	return &Repo{db: db}
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error) {
	const query = `SELECT * FROM organization WHERE id = $1`

	var mdl model

	row := r.db.QueryRow(ctx, query, id)
	if err := row.Scan(&mdl); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return mdl.toOrganization(), nil
}
