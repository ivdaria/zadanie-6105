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

func (s *Server) SubmitBidFeedback(ctx echo.Context, bidId api.BidId, params api.SubmitBidFeedbackParams) error {
	rctx := ctx.Request().Context()

	// есть ли пользователь с таким именем
	employee, err := s.employees.GetByUserName(rctx, params.Username)
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
	if bid.Status != entity.BidStatusPublished && bid.Status != entity.BidStatusCanceled {
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: "can not see bid because tender organization can see bids only if they are published or cancelled",
		})
	}

	// Проверка существования тендера
	tender, err := s.tenders.GetTenderByID(rctx, bid.TenderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Reason: fmt.Sprintf("no tender with this ID: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("get tender by ID: %v", err.Error()),
		})
	}

	// проверяем, является ли пользователь ответственным в организации, которая опубликовала тендер
	if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, employee.ID, tender.OrganizationID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
				Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
		})

	}

	if err := s.feedbacks.CreateFeedback(rctx, &entity.Feedback{
		BidID:            bid.ID,
		FeedbackAuthorID: employee.ID,
		Comment:          params.BidFeedback,
	}); err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("create feedback: %v", err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, api.Bid{
		Id:          bid.ID.String(),
		Name:        bid.Name,
		Description: bid.Description,
		TenderId:    bid.TenderID.String(),
		AuthorId:    bid.CreatorID.String(),
		AuthorType:  api.BidAuthorType(bid.AuthorType),
		Status:      api.BidStatus(bid.Status),
		Version:     api.BidVersion(bid.Version),
		CreatedAt:   bid.CreatedAt.Format(time.RFC3339),
	})
}
