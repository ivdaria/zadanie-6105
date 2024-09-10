package gateway

import (
	"context"
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/pkg/api"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type tenderRepo interface {
	CreateTender(ctx context.Context, tender *entity.Tender) (*entity.Tender, error)
	GetAllTenders(ctx context.Context, filter entity.GetTendersFilter, pagination entity.Pagination) ([]*entity.Tender, error)
	GetTendersByUsername(ctx context.Context, username string, pagination entity.Pagination) ([]*entity.Tender, error)
	GetStatusByID(ctx context.Context, id uuid.UUID) (entity.TenderStatus, error)
}

type employeeRepo interface {
	GetByUserName(ctx context.Context, username string) (*entity.Employee, error)
}

type organizationRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error)
}

type bidRepo interface {
	CreateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error)
}

type Server struct {
	tenders       tenderRepo
	employees     employeeRepo
	organizations organizationRepo
	bids          bidRepo
}

func NewServer(
	tenders tenderRepo,
	employees employeeRepo,
	organizations organizationRepo,
	bids bidRepo,
) *Server {
	return &Server{
		tenders:       tenders,
		employees:     employees,
		organizations: organizations,
		bids:          bids,
	}
}

func (s *Server) GetUserBids(ctx echo.Context, params api.GetUserBidsParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) CreateBid(ctx echo.Context) error {
	//rctx := ctx.Request().Context()
	//var body api.CreateBidJSONBody
	//if err := ctx.Bind(&body); err != nil {
	//	return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
	//		Reason: fmt.Sprintf("failed to bind body: %v", err.Error()),
	//	})
	//}
	//
	//employee, err := s.employees.GetByUserName(rctx, body.CreatorUsername)
	//if err != nil {
	//	return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
	//		Reason: fmt.Sprintf("failed to get employee by username: %v", err.Error()),
	//	})
	//}
	//
	//// TODO сделать маппинг отдельной функцией CreateTenderJSONBody->entity
	//bid := &entity.Bid{
	//	ID:             0,
	//	TenderID:       0,
	//	CreatorID:      0,
	//	OrganizationID: 0,
	//	Decision:       "",
	//	Status:         "",
	//	AuthorType:     "",
	//	Version:        1,
	//}
	//
	//bid.ID = s.bids.CreateBid(bid)
	//
	//return ctx.JSON(http.StatusOK, api.Bid{
	//	AuthorId:    "",
	//	AuthorType:  "",
	//	CreatedAt:   "",
	//	Description: "",
	//	Id:          "",
	//	Name:        "",
	//	Status:      "",
	//	TenderId:    "",
	//	Version:     1,
	//})
	return nil
}

func (s *Server) EditBid(ctx echo.Context, bidId api.BidId, params api.EditBidParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) SubmitBidFeedback(ctx echo.Context, bidId api.BidId, params api.SubmitBidFeedbackParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) RollbackBid(ctx echo.Context, bidId api.BidId, version int32, params api.RollbackBidParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetBidStatus(ctx echo.Context, bidId api.BidId, params api.GetBidStatusParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) UpdateBidStatus(ctx echo.Context, bidId api.BidId, params api.UpdateBidStatusParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) SubmitBidDecision(ctx echo.Context, bidId api.BidId, params api.SubmitBidDecisionParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetBidsForTender(ctx echo.Context, tenderId api.TenderId, params api.GetBidsForTenderParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetBidReviews(ctx echo.Context, tenderId api.TenderId, params api.GetBidReviewsParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) CheckServer(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, nil)
}

func (s *Server) GetTenders(ctx echo.Context, params api.GetTendersParams) error {
	rctx := ctx.Request().Context()

	var serviceTypesFilter entity.GetTendersFilter
	if params.ServiceType != nil {
		var serviceTypes []string
		for _, serviceType := range *params.ServiceType {
			serviceTypes = append(serviceTypes, string(serviceType))
		}
		serviceTypesFilter.ServiceTypes = &serviceTypes
	}

	allTenders, err := s.tenders.GetAllTenders(
		rctx,
		serviceTypesFilter,
		entity.NewPagination(params.Limit, params.Offset),
	)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get all tenders: %v", err.Error()),
		})
	}

	if len(allTenders) == 0 {
		return ctx.JSON(http.StatusOK, []interface{}{})
	}

	return ctx.JSON(http.StatusOK, allTenders)
}

func (s *Server) GetUserTenders(ctx echo.Context, params api.GetUserTendersParams) error {
	rctx := ctx.Request().Context()

	if params.Username == nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{})
	}

	//TODO проверить, существует ли пользователь

	tenders, err := s.tenders.GetTendersByUsername(rctx, *params.Username, entity.NewPagination(params.Limit, params.Offset))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get user's tenders: %v", err.Error()),
		})
	}

	if len(tenders) == 0 {
		return ctx.JSON(http.StatusOK, []interface{}{})
	}

	return ctx.JSON(http.StatusOK, tenders)
}

func (s *Server) CreateTender(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	var body api.CreateTenderJSONBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to bind body: %v", err.Error()),
		})
	}

	organizationID, err := uuid.Parse(body.OrganizationId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse OrganizationId: %v", err.Error()),
		})
	}

	organization, err := s.organizations.GetByID(rctx, organizationID)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get organization by id: %v", err.Error()),
		})
	}

	employee, err := s.employees.GetByUserName(rctx, body.CreatorUsername)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get employee by username: %v", err.Error()),
		})
	}

	// TODO сделать маппинг отдельной функцией CreateTenderJSONBody->entity
	tender := &entity.Tender{
		Name:           body.Name,
		Description:    body.Description,
		ServiceType:    entity.ServiceType(body.ServiceType),
		Status:         entity.TenderStatus(body.Status),
		OrganizationID: organization.ID,
		CreatorID:      employee.ID,
		Version:        1,
	}

	tender, err = s.tenders.CreateTender(rctx, tender)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to create tender: %v", err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, api.Tender{
		CreatedAt:      tender.CreatedAt.Format(time.RFC3339),
		Description:    tender.Description,
		Id:             tender.ID.String(),
		Name:           body.Name,
		OrganizationId: organization.ID.String(),
		ServiceType:    api.TenderServiceType(tender.ServiceType),
		Status:         api.TenderStatus(tender.Status),
		Version:        api.TenderVersion(tender.Version),
	})
}

func (s *Server) EditTender(ctx echo.Context, tenderId api.TenderId, params api.EditTenderParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) RollbackTender(ctx echo.Context, tenderId api.TenderId, version int32, params api.RollbackTenderParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetTenderStatus(ctx echo.Context, tenderId api.TenderId, params api.GetTenderStatusParams) error {
	rctx := ctx.Request().Context()

	var body api.GetTenderStatusParams
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to bind body: %v", err.Error()),
		})
	}

	tenderIDParsed, err := uuid.Parse(tenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse tenderID: %v", err.Error()),
		})
	}

	tendersStatus, err := s.tenders.GetStatusByID(rctx, tenderIDParsed)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender's status by tender ID: %v", err.Error()),
		})
	}
	return ctx.JSON(http.StatusOK, tendersStatus)
}

func (s *Server) UpdateTenderStatus(ctx echo.Context, tenderId api.TenderId, params api.UpdateTenderStatusParams) error {
	//TODO implement me
	panic("implement me")
}
