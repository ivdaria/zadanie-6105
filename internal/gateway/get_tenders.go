package gateway

import (
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/pkg/api"
	"github.com/labstack/echo/v4"
	"net/http"
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

	if len(allTenders) == 0 {
		return ctx.JSON(http.StatusOK, []interface{}{})
	}

	return ctx.JSON(http.StatusOK, allTenders)
}
