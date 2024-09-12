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
	"time"
)

func (s *Server) EditTender(ctx echo.Context, tenderId api.TenderId, params api.EditTenderParams) error {
	rctx := ctx.Request().Context()
	var body api.EditTenderJSONBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to bind body: %v", err.Error()),
		})
	}

	//проверить, существует ли пользователь

	employee, err := s.employees.GetByUserName(rctx, params.Username)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Reason: fmt.Sprintf("no employee with: %v", err.Error()),
		})
	}

	// существует ли тендер
	oldTenderID, err := uuid.Parse(tenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse tender ID: %v", err.Error()),
		})
	}

	oldTender, err := s.tenders.GetTenderByID(rctx, oldTenderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Reason: fmt.Sprintf("no tender with this ID: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get tender by id: %v", err.Error()),
		})

	}

	//есть ли права у пользователя

	if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, employee.ID, oldTender.OrganizationID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not an organization's responsible: %v", err.Error()),
		})
	}

	patchedTender := oldTender.Patch(body.Name, body.Description, (*entity.ServiceType)(body.ServiceType))
	patchedTender, err = s.tenders.UpdateTender(rctx, patchedTender)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to update tender: %v", err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, api.Tender{
		CreatedAt:      patchedTender.CreatedAt.Format(time.RFC3339),
		Description:    patchedTender.Description,
		Id:             patchedTender.ID.String(),
		Name:           patchedTender.Name,
		OrganizationId: patchedTender.ID.String(),
		ServiceType:    api.TenderServiceType(patchedTender.ServiceType),
		Status:         api.TenderStatus(patchedTender.Status),
		Version:        api.TenderVersion(patchedTender.Version),
	})
}
