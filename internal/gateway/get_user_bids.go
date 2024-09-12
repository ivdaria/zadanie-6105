package gateway

import (
	"errors"
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/pkg/api"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) GetUserBids(ctx echo.Context, params api.GetUserBidsParams) error {
	rctx := ctx.Request().Context()

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

	bids, err := s.bids.GetBidsByUsername(rctx, *params.Username, entity.NewPagination(params.Limit, params.Offset))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get user's bids: %v", err.Error()),
		})
	}

	if len(bids) == 0 {
		return ctx.JSON(http.StatusOK, []interface{}{})
	}

	return ctx.JSON(http.StatusOK, bids)
}
