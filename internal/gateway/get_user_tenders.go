package gateway

import (
	"errors"
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/pkg/api"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (s *Server) GetUserTenders(ctx echo.Context, params api.GetUserTendersParams) error {
	rctx := ctx.Request().Context()

	if params.Username == nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("add username"),
		})
	}

	_, err := s.employees.GetByUserName(rctx, *params.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Reason: fmt.Sprintf("no employee with: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get employee: %v", err.Error()),
		})
	}

	tenders, err := s.tenders.GetTendersByUsername(rctx, *params.Username, entity.NewPagination(params.Limit, params.Offset))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get user's tenders: %v", err.Error()),
		})
	}

	apiTenders := make([]api.Tender, 0, len(tenders))
	for _, tender := range tenders {
		apiTenders = append(apiTenders, api.Tender{
			CreatedAt:      tender.CreatedAt.Format(time.RFC3339),
			Description:    tender.Description,
			Id:             tender.ID.String(),
			Name:           tender.Name,
			OrganizationId: tender.ID.String(),
			ServiceType:    api.TenderServiceType(tender.ServiceType),
			Status:         api.TenderStatus(tender.Status),
			Version:        api.TenderVersion(tender.Version),
		})
	}

	return ctx.JSON(http.StatusOK, apiTenders)
}
