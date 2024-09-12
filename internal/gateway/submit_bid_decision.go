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

func (s *Server) SubmitBidDecision(ctx echo.Context, bidId api.BidId, params api.SubmitBidDecisionParams) error {
	rctx := ctx.Request().Context()

	// проверка наличия Username
	if params.Username == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("add username"),
		})
	}

	decisionSubmitUser, err := s.employees.GetByUserName(rctx, params.Username)
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

	bidSubmitID, err := uuid.Parse(bidId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse bid ID: %v", err.Error()),
		})
	}

	bidForSubmition, err := s.bids.GetBidByID(rctx, bidSubmitID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("no bid with this ID: %v", err.Error()),
		})
	}

	tenderForBid, err := s.tenders.GetTenderByID(rctx, bidForSubmition.TenderID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender: %v", err.Error()),
		})
	}

	if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, decisionSubmitUser.ID, tenderForBid.OrganizationID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("check is user organization responsible: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not organization's responsible or not an author of bid"),
		})
	}

	err = s.bids.UpdateBidDecision(rctx, bidForSubmition.ID, entity.BidDecision(params.Decision))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to update bid's decision: %v", err.Error()),
		})
	}

	updatedBid, err := s.bids.GetBidByID(rctx, bidForSubmition.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get updated bid: %v", err.Error()),
		})
	}

	if updatedBid.Decision == entity.BidDecisionApproved {
		err = s.tenders.UpdateTenderStatus(rctx, tenderForBid.ID, entity.TenderStatusClosed)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to update tender's status: %v", err.Error()),
			})
		}

		_, err = s.tenders.GetTenderByID(rctx, tenderForBid.ID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get updated tender: %v", err.Error()),
			})
		}
	}

	return ctx.JSON(http.StatusOK, api.Bid{
		AuthorType:  api.BidAuthorType(updatedBid.AuthorType),
		AuthorId:    updatedBid.CreatorID.String(),
		TenderId:    updatedBid.TenderID.String(),
		CreatedAt:   updatedBid.CreatedAt.Format(time.RFC3339),
		Description: updatedBid.Description,
		Id:          updatedBid.ID.String(),
		Name:        updatedBid.Name,
		Status:      api.BidStatus(updatedBid.Status),
		Version:     api.BidVersion(updatedBid.Version),
	})
}
