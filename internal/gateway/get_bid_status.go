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
)

func (s *Server) GetBidStatus(ctx echo.Context, bidId api.BidId, params api.GetBidStatusParams) error {
	rctx := ctx.Request().Context()

	// проверка наличия Username
	if params.Username == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("add username"),
		})
	}

	// есть ли пользователь с таким именем
	userToGetStatus, err := s.employees.GetByUserName(rctx, params.Username)
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
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("no bid with this ID: %v", err.Error()),
		})
	}

	// для проверки статуса предложения пользователь должен быть либо автором предложения (если AuthorType - User)
	// либо ответственным в организации (если AuthorType - Organization)
	// либо ответственным в организации, которая разместила тендер, связанный с предложением

	var (
		onBidSide    bool = true
		onOrgsSide   bool = true
		onTenderSide bool = true
	)

	if bid.AuthorType == entity.BidAuthorTypeUser {
		if bid.CreatorID != userToGetStatus.ID {
			onBidSide = false
		}
	} else {
		// проверяем, является ли пользователь ответственным организации
		oldBitOrganization, err := s.organizationResponsibles.GetOrganizationResponsibleByUserID(rctx, bid.CreatorID)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get organization ID: %v", err.Error()),
			})
		}

		if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, userToGetStatus.ID, oldBitOrganization.OrganizationID); errors.Is(err, pgx.ErrNoRows) {
			if !errors.Is(err, pgx.ErrNoRows) {
				return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
					Reason: fmt.Sprintf("check is user organization responsible: %v", err.Error()),
				})
			}
			onBidSide = false
		}
	}

	if err := s.organizationResponsibles.IsUserResponsible(rctx, userToGetStatus.ID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("check is user responsible: %v", err.Error()),
			})
		}
		onBidSide = false
	}

	tenderForBid, err := s.tenders.GetTenderByID(rctx, bid.TenderID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender: %v", err.Error()),
		})
	}

	if !onBidSide || !onOrgsSide {
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not organization's responsible or not an author of bid"),
		})
	}

	if bid.Status != entity.BidStatusPublished {
		if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, userToGetStatus.ID, tenderForBid.OrganizationID); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
					Reason: fmt.Sprintf("check is user responsible: %v", err.Error()),
				})
			}
			onTenderSide = false
		}

		if !onBidSide || !onTenderSide {
			return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
				Reason: fmt.Sprintf("user is not organization's responsible or not an author of bid"),
			})
		}

	}

	return ctx.JSON(http.StatusOK, bid.Status)
}
