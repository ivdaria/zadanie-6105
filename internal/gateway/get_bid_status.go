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
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Reason: fmt.Sprintf("no bid with this ID: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get bid: %v", err.Error()),
		})
	}

	tenderForBid, err := s.tenders.GetTenderByID(rctx, bid.TenderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Reason: fmt.Sprintf("no teder with this ID: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender: %v", err.Error()),
		})
	}

	// для проверки статуса предложения пользователь должен быть либо автором предложения (если AuthorType - User)
	// либо ответственным в организации (если AuthorType - Organization)
	// либо ответственным в организации, которая разместила тендер, связанный с предложением
	isBidAuthor := true
	if bid.AuthorType == entity.BidAuthorTypeUser {
		if bid.CreatorID != userToGetStatus.ID {
			isBidAuthor = false
		}
	} else {
		// проверяем, является ли пользователь ответственным организации
		oldBidOrganization, err := s.organizationResponsibles.GetOrganizationResponsibleByUserID(rctx, bid.CreatorID)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get organization ID: %v", err.Error()),
			})
		}

		if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, userToGetStatus.ID, oldBidOrganization.OrganizationID); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
					Reason: fmt.Sprintf("check is user organization responsible: %v", err.Error()),
				})
			}
			isBidAuthor = false
		}
	}

	// Если не автор или организация, создавшая бид, то проверяем, является ли пользователь ответственным за тендер
	if !isBidAuthor {
		if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, userToGetStatus.ID, tenderForBid.OrganizationID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
					Reason: fmt.Sprintf("can not see bid because it's not a user in tender organization or bid author or bid author org: %v", err.Error()),
				})
			}
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("check is user organization responsible: %v", err.Error()),
			})
		}

		if bid.Status != entity.BidStatusPublished && bid.Status != entity.BidStatusCanceled {
			return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
				Reason: "can not see bid because tender organization can see bids only if they are published or cancelled",
			})
		}
	}

	return ctx.JSON(http.StatusOK, bid.Status)
}
