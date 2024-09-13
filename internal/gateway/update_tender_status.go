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

func (s *Server) UpdateTenderStatus(ctx echo.Context, tenderId api.TenderId, params api.UpdateTenderStatusParams) error {
	rctx := ctx.Request().Context()

	if params.Username == "" || params.Status == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: "add status or username",
		})
	}

	employee, err := s.employees.GetByUserName(rctx, params.Username)
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

	tenderUUID, err := uuid.Parse(tenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse tenderID: %v", err.Error()),
		})
	}

	tender, err := s.tenders.GetTenderByID(rctx, tenderUUID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender, please, check ID: %v", err.Error()),
		})
	}

	if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, employee.ID, tender.OrganizationID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
		})
	}

	err = s.tenders.UpdateTenderStatus(rctx, tender.ID, entity.TenderStatus(params.Status))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to update tender's status: %v", err.Error()),
		})
	}

	updatedTender, err := s.tenders.GetTenderByID(rctx, tenderUUID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get updated tender: %v", err.Error()),
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
