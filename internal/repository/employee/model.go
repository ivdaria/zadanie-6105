package employee

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type model struct {
	ID        uuid.UUID
	UserName  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type models []*model

func (m *model) ScanRow(rows pgx.Rows) error {
	return rows.Scan(&m.ID, &m.UserName, &m.FirstName, &m.LastName, &m.CreatedAt, &m.UpdatedAt)
}

func (m *model) toEmployee() *entity.Employee {
	return &entity.Employee{
		ID:        m.ID,
		UserName:  m.UserName,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (mdls models) toEmployees() []*entity.Employee {
	if len(mdls) == 0 {
		return nil
	}

	result := make([]*entity.Employee, 0, len(mdls))
	for _, m := range mdls {
		result = append(result, m.toEmployee())
	}

	return result
}

func modelFromEmployee(item *entity.Employee) *model {
	return &model{
		ID:        item.ID,
		UserName:  item.UserName,
		FirstName: item.FirstName,
		LastName:  item.LastName,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
