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

func (s *Server) CreateBid(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	var body api.CreateBidJSONBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to bind body: %v", err.Error()),
		})
	}

	tenderID, err := uuid.Parse(body.TenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse TenderID: %v", err.Error()),
		})
	}

	// Проверка существования тендера
	tender, err := s.tenders.GetTenderByID(rctx, tenderID)
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
	// Проверка, что тендер опубликован
	if tender.Status != entity.TenderStatusPublished {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("tender is not published"),
		})
	}

	creatorID, err := uuid.Parse(body.AuthorId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse CreatorID: %v", err.Error()),
		})
	}

	_, err = s.employees.GetEmployeeByID(rctx, creatorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Reason: fmt.Sprintf("user is not an organization's responsible"),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get employee by id: %v", err.Error()),
		})
	}

	if err := s.organizationResponsibles.IsUserResponsible(rctx, creatorID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
				Reason: fmt.Sprintf("user is not an organization's responsible"),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
		})
	}

	bid := &entity.Bid{
		Name:        body.Name,
		TenderID:    tenderID,
		CreatorID:   creatorID,
		Description: body.Description,
		Decision:    entity.BidDecisionNone,
		Status:      entity.BidStatusCreated,
		AuthorType:  entity.BidAuthorType(body.AuthorType),
		Version:     1,
	}

	bid, err = s.bids.CreateBid(rctx, bid)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to create bid: %v", err.Error()),
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
