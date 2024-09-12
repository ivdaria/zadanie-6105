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

func (s *Server) EditBid(ctx echo.Context, bidId api.BidId, params api.EditBidParams) error {
	rctx := ctx.Request().Context()

	var body api.EditBidJSONBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to bind body: %v", err.Error()),
		})
	}

	// есть ли пользователь с таким именем
	bidChanger, err := s.employees.GetByUserName(rctx, params.Username)
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

	// существует ли предложение
	oldBidID, err := uuid.Parse(bidId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse bid ID: %v", err.Error()),
		})
	}

	oldBid, err := s.bids.GetBidByID(rctx, oldBidID)
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

	patchedBid := oldBid.Patch(body.Name, body.Description)
	patchedBid, err = s.bids.UpdateBid(rctx, patchedBid)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to update bid: %v", err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, api.Bid{
		Id:          patchedBid.ID.String(),
		Name:        patchedBid.Name,
		Description: patchedBid.Description,
		TenderId:    patchedBid.TenderID.String(),
		AuthorId:    patchedBid.CreatorID.String(),
		AuthorType:  api.BidAuthorType(patchedBid.AuthorType),
		Status:      api.BidStatus(patchedBid.Status),
		Version:     api.BidVersion(patchedBid.Version),
		CreatedAt:   patchedBid.CreatedAt.Format(time.RFC3339),
	})
}
