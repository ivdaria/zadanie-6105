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

func (s *Server) GetBidsForTender(ctx echo.Context, tenderId api.TenderId, params api.GetBidsForTenderParams) error {
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

	// парсим tenderId на uuid, получаем тендер, проверяем, есть ли такой тендер вообще
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
				Reason: fmt.Sprintf("no tender with this ID: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get tender by id: %v", err.Error()),
		})
	}

	// проверяем, является ли пользователь ответственным в организации, которая опубликовала тендер
	if err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, employee.ID, tender.OrganizationID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
				Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
		})
	}

	// получаем предложения для тендера
	bids, err := s.bids.GetBidsByTenderID(rctx, tenderIDParsed, entity.NewPagination(params.Limit, params.Offset))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get bids for tender: %v", err.Error()),
		})
	}

	apiBids := make([]api.Bid, 0, len(bids))
	for _, bid := range bids {
		apiBids = append(apiBids, api.Bid{
			AuthorId:    bid.CreatorID.String(),
			AuthorType:  api.BidAuthorType(bid.AuthorType),
			CreatedAt:   bid.CreatedAt.Format(time.RFC3339),
			Description: bid.Description,
			Id:          bid.ID.String(),
			Name:        bid.Name,
			Status:      api.BidStatus(bid.Status),
			TenderId:    bid.TenderID.String(),
			Version:     api.BidVersion(bid.Version),
		})
	}

	return ctx.JSON(http.StatusOK, apiBids)

}
