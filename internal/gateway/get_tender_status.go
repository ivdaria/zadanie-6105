package gateway

import (
	"errors"
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/pkg/api"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) GetTenderStatus(ctx echo.Context, tenderId api.TenderId, params api.GetTenderStatusParams) error {
	rctx := ctx.Request().Context()

	var body api.GetTenderStatusParams
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to bind body: %v", err.Error()),
		})
	}

	// поиск пользователя по имени - если нет, то 401
	tenderCreator, err := s.employees.GetByUserName(rctx, *params.Username)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get employee by username: %v", err.Error()),
		})
	}

	//парсим tenderId, смотрим, есть ли такой тендер
	tenderIDParsed, err := uuid.Parse(tenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse tenderID: %v", err.Error()),
		})
	}

	tender, err := s.tenders.GetTenderByID(rctx, tenderIDParsed)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender's status by tender ID: %v", err.Error()),
		})
	}

	if tender.Status != entity.TenderStatusPublished {
		// если ID пользователя - не равно ID автора тендера и не в списке ответственных за организацию, то 403
		if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, tenderCreator.ID, tender.OrganizationID); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
					Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
				})
			}
			return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
				Reason: fmt.Sprintf("user is not an organization's responsible: %v", err.Error()),
			})
		}
	}

	return ctx.JSON(http.StatusOK, tender.Status)
}
