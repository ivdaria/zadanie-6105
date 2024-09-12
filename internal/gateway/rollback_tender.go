package gateway

import (
	"errors"
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/pkg/api"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (s *Server) RollbackTender(ctx echo.Context, tenderId api.TenderId, version int32, params api.RollbackTenderParams) error {
	rctx := ctx.Request().Context()

	if params.Username == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("add username"),
		})
	}
	if version < 1 {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("invalid version. Version should be >= 1"),
		})
	}

	tenderChanger, err := s.employees.GetByUserName(rctx, params.Username)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get employee by username: %v", err.Error()),
		})
	}

	tenderIDParsed, err := uuid.Parse(tenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse tenderID: %v", err.Error()),
		})
	}

	tender, err := s.tenders.GetTenderByID(rctx, tenderIDParsed)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender, please, check ID: %v", err.Error()),
		})
	}

	if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, tenderChanger.ID, tender.OrganizationID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not an organization's responsible: %v", err.Error()),
		})
	}

	tenderToRollback, err := s.tenders.GetTenderByIDAndVersion(rctx, tenderIDParsed, int(version))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender's version, please, check parametrs: %v", err.Error()),
		})
	}

	tenderToRollback = tender.Rollback(tenderToRollback)

	updatedTender, err := s.tenders.UpdateTender(rctx, tenderToRollback)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to update tender: %v", err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, api.Tender{
		CreatedAt:      updatedTender.CreatedAt.Format(time.RFC3339),
		Description:    updatedTender.Description,
		Id:             updatedTender.ID.String(),
		Name:           updatedTender.Name,
		OrganizationId: updatedTender.ID.String(),
		ServiceType:    api.TenderServiceType(updatedTender.ServiceType),
		Status:         api.TenderStatus(updatedTender.Status),
		Version:        api.TenderVersion(updatedTender.Version),
	})
}
