package employee

import (
	"context"
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/jackc/pgx/v5"
)

// Repo содержит в себе логику хранения сотрудников
type Repo struct {
	db *pgx.Conn
}

func NewRepo(db *pgx.Conn) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetByUserName(ctx context.Context, username string) (*entity.Employee, error) {
	const query = `SELECT * FROM employee WHERE username = $1`

	var mdl model

	row := r.db.QueryRow(ctx, query, username)
	if err := row.Scan(&mdl); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return mdl.toEmployee(), nil
}
