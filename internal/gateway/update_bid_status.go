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

func (s *Server) UpdateBidStatus(ctx echo.Context, bidId api.BidId, params api.UpdateBidStatusParams) error {
	rctx := ctx.Request().Context()

	// есть ли пользователь с таким именем
	user, err := s.employees.GetByUserName(rctx, params.Username)
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

	// есть ли такое предложение
	bidIDParsed, err := uuid.Parse(bidId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse bid ID: %v", err.Error()),
		})
	}

	bid, err := s.bids.GetBidByID(rctx, bidIDParsed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Reason: fmt.Sprintf("no bid with this ID: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get bid: %v", err.Error()),
		})
	}

	if bid.AuthorType == entity.BidAuthorTypeUser {
		if bid.CreatorID != user.ID {
			return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
				Reason: "you are not bid author",
			})
		}
	} else {
		// проверяем, является ли пользователь ответственным организации
		oldBidOrganization, err := s.organizationResponsibles.GetOrganizationResponsibleByUserID(rctx, bid.CreatorID)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get organization ID: %v", err.Error()),
			})
		}

		if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, user.ID, oldBidOrganization.OrganizationID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
					Reason: fmt.Sprintf("you are not responsible for bid: %v", err.Error()),
				})
			}
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("check is user organization responsible: %v", err.Error()),
			})
		}
	}

	if err := s.bids.UpdateBidStatus(rctx, bid.ID, entity.BidStatus(params.Status)); err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("update bid status: %v", err.Error()),
		})
	}

	updatedBid, err := s.bids.GetBidByID(rctx, bidIDParsed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Reason: fmt.Sprintf("no bid with this ID: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get updated bid: %v", err.Error()),
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
