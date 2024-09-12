package gateway

import (
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/pkg/api"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (s *Server) GetTenders(ctx echo.Context, params api.GetTendersParams) error {
	rctx := ctx.Request().Context()

	var serviceTypesFilter entity.GetTendersFilter
	if params.ServiceType != nil {
		var serviceTypes []string
		for _, serviceType := range *params.ServiceType {
			serviceTypes = append(serviceTypes, string(serviceType))
		}
		serviceTypesFilter.ServiceTypes = &serviceTypes
	}

	allTenders, err := s.tenders.GetAllTenders(
		rctx,
		serviceTypesFilter,
		entity.NewPagination(params.Limit, params.Offset),
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get all tenders: %v", err.Error()),
		})
	}

	apiTenders := make([]api.Tender, 0, len(allTenders))
	for _, tender := range allTenders {
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
