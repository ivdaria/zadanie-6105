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

func (s *Server) GetBidReviews(ctx echo.Context, tenderId api.TenderId, params api.GetBidReviewsParams) error {
	rctx := ctx.Request().Context()

	//проверить, существует ли запрашивающий пользователь
	requesterEmployee, err := s.employees.GetByUserName(rctx, params.RequesterUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Reason: fmt.Sprintf("no requesterEmployee with: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get requesterEmployee by id: %v", err.Error()),
		})
	}

	//проверить, существует ли автор фидбэка
	authorEmployee, err := s.employees.GetByUserName(rctx, params.AuthorUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Reason: fmt.Sprintf("no authorEmployee with: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get authorEmployee by id: %v", err.Error()),
		})
	}

	// существует ли тендер
	tenderIDParsed, err := uuid.Parse(tenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse tender ID: %v", err.Error()),
		})
	}

	tender, err := s.tenders.GetTenderByID(rctx, tenderIDParsed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get tender by tender ID: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get tender by id: %v", err.Error()),
		})
	}

	// если ID пользователя - не в списке ответственных за организацию, то 403
	if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, requesterEmployee.ID, tender.OrganizationID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
				Reason: fmt.Sprintf("user is not an organization's responsible: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
		})
	}

	feedbacks, err := s.feedbacks.GetFeedbackByTenderIDAndAuthor(
		rctx,
		tender.ID,
		authorEmployee.ID,
		entity.NewPagination(params.Limit, params.Offset),
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get feedbacks: %v", err.Error()),
		})
	}

	apiBidReviews := make([]api.BidReview, 0, len(feedbacks))
	for _, feedback := range feedbacks {
		apiBidReviews = append(apiBidReviews, api.BidReview{
			CreatedAt:   feedback.CreatedAt.Format(time.RFC3339),
			Description: feedback.Comment,
			Id:          feedback.ID.String(),
		})
	}

	return ctx.JSON(http.StatusOK, apiBidReviews)
}
