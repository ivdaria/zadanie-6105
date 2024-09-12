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

func (s *Server) RollbackBid(ctx echo.Context, bidId api.BidId, version int32, params api.RollbackBidParams) error {
	rctx := ctx.Request().Context()

	if version < 1 {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("invalid version. Version should be >= 1"),
		})
	}

	bidChanger, err := s.employees.GetByUserName(rctx, params.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get employee by username: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get empoloyee by username: %v", err.Error()),
		})
	}

	bidIDParsed, err := uuid.Parse(bidId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse bid ID: %v", err.Error()),
		})
	}

	oldBid, err := s.bids.GetBidByID(rctx, bidIDParsed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Reason: fmt.Sprintf("no bid with this ID: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get bid by id: %v", err.Error()),
		})
	}

	canEdit, err := s.userCanEditBidChecker.IsUserCanEditBid(rctx, oldBid, bidChanger)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("check is user can edit bid: %v", err.Error()),
		})
	}
	if !canEdit {
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: "can not edit bid",
		})
	}

	bidToRollback, err := s.bids.GetBidByIDAndVersion(rctx, bidIDParsed, int(version))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get bid's version, please, check parametrs: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get bid by id and version: %v", err.Error()),
		})
	}

	bidToRollback = oldBid.Rollback(bidToRollback)
	updatedBid, err := s.bids.UpdateBid(rctx, bidToRollback)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to update bid: %v", err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, api.Bid{
		Id:          updatedBid.ID.String(),
		Name:        updatedBid.Name,
		Description: updatedBid.Description,
		TenderId:    updatedBid.TenderID.String(),
		AuthorId:    updatedBid.CreatorID.String(),
		AuthorType:  api.BidAuthorType(updatedBid.AuthorType),
		Status:      api.BidStatus(updatedBid.Status),
		Version:     api.BidVersion(updatedBid.Version),
		CreatedAt:   updatedBid.CreatedAt.Format(time.RFC3339),
	})
}
