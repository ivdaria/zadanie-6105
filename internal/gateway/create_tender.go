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

func (s *Server) CreateTender(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	var body api.CreateTenderJSONBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to bind body: %v", err.Error()),
		})
	}

	organizationID, err := uuid.Parse(body.OrganizationId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse OrganizationId: %v", err.Error()),
		})
	}

	organization, err := s.organizations.GetByID(rctx, organizationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get organization by id: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get organization: %v", err.Error()),
		})
	}

	employee, err := s.employees.GetByUserName(rctx, body.CreatorUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get employee by username: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get employee by username: %v", err.Error()),
		})
	}

	// доступно только ответственным за организацию

	if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, employee.ID, organizationID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not an organization's responsible"),
		})
	}

	// TODO сделать маппинг отдельной функцией CreateTenderJSONBody->entity
	tender := &entity.Tender{
		Name:           body.Name,
		Description:    body.Description,
		ServiceType:    entity.ServiceType(body.ServiceType),
		Status:         entity.TenderStatusCreated,
		OrganizationID: organization.ID,
		CreatorID:      employee.ID,
		Version:        1,
	}

	tender, err = s.tenders.CreateTender(rctx, tender)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to create tender: %v", err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, api.Tender{
		CreatedAt:      tender.CreatedAt.Format(time.RFC3339),
		Description:    tender.Description,
		Id:             tender.ID.String(),
		Name:           body.Name,
		OrganizationId: organization.ID.String(),
		ServiceType:    api.TenderServiceType(tender.ServiceType),
		Status:         api.TenderStatus(tender.Status),
		Version:        api.TenderVersion(tender.Version),
	})
}
