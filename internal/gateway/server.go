package gateway

import (
	"context"
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

type tenderRepo interface {
	CreateTender(ctx context.Context, tender *entity.Tender) (*entity.Tender, error)
	GetAllTenders(ctx context.Context, filter entity.GetTendersFilter, pagination entity.Pagination) ([]*entity.Tender, error)
	GetTendersByUsername(ctx context.Context, username string, pagination entity.Pagination) ([]*entity.Tender, error)
	GetTenderByID(ctx context.Context, id uuid.UUID) (*entity.Tender, error)
	UpdateTender(ctx context.Context, tender *entity.Tender) (*entity.Tender, error)
	UpdateTenderStatus(ctx context.Context, id uuid.UUID, newStatus entity.TenderStatus) error
	GetTenderByIDAndVersion(ctx context.Context, id uuid.UUID, version int) (*entity.Tender, error)
}

type employeeRepo interface {
	GetByUserName(ctx context.Context, username string) (*entity.Employee, error)
	GetEmployeeByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error)
}

type organizationRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error)
}

type bidRepo interface {
	CreateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error)
}

type organizationResponsibleRepo interface {
	IsUserOrganizationResponsible(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (bool, error)
}

type Server struct {
	tenders                  tenderRepo
	employees                employeeRepo
	organizations            organizationRepo
	bids                     bidRepo
	organizationResponsibles organizationResponsibleRepo
}

func NewServer(
	tenders tenderRepo,
	employees employeeRepo,
	organizations organizationRepo,
	bids bidRepo,
	organizationResponsibles organizationResponsibleRepo,
) *Server {
	return &Server{
		tenders:                  tenders,
		employees:                employees,
		organizations:            organizations,
		bids:                     bids,
		organizationResponsibles: organizationResponsibles,
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
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("add username"),
		})
	}

	_, err := s.employees.GetByUserName(rctx, *params.Username)
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

	tenders, err := s.tenders.GetTendersByUsername(rctx, *params.Username, entity.NewPagination(params.Limit, params.Offset))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
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
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get organization by id: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get organization: %v", err.Error()),
		})
	}

	employee, err := s.employees.GetByUserName(rctx, body.CreatorUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Reason: fmt.Sprintf("failed to get employee by username: %v", err.Error()),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("get employee by username: %v", err.Error()),
		})
	}

	//TODO доступно только ответственным за организацию

	isResponsible, err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, employee.ID, organizationID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("check is responsible: %v", err.Error()),
		})
	}
	if !isResponsible {
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not an organization's responsible"),
		})
	}

	// TODO сделать маппинг отдельной функцией CreateTenderJSONBody->entity
	tender := &entity.Tender{
		Name:           body.Name,
		Description:    body.Description,
		ServiceType:    entity.ServiceType(body.ServiceType),
		Status:         entity.TenderStatusCreated,
		OrganizationID: organization.ID,
		CreatorID:      employee.ID,
		Version:        1,
	}

	tender, err = s.tenders.CreateTender(rctx, tender)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
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
	rctx := ctx.Request().Context()
	var body api.EditTenderJSONBody
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to bind body: %v", err.Error()),
		})
	}

	//проверить, существует ли пользователь

	employee, err := s.employees.GetByUserName(rctx, params.Username)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Reason: fmt.Sprintf("no employee with: %v", err.Error()),
		})
	}

	// существует ли тендер
	oldTenderID, err := uuid.Parse(tenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse tender ID: %v", err.Error()),
		})
	}

	oldTender, err := s.tenders.GetTenderByID(rctx, oldTenderID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("no tender with this ID: %v", err.Error()),
		})
	}

	//есть ли права у пользователя

	isResponsible, err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, employee.ID, oldTender.OrganizationID)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not an organization's responsible: %v", err.Error()),
		})
	}
	if !isResponsible {
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not an organization's responsible"),
		})
	}

	patchedTender := oldTender.Patch(body.Name, body.Description, (*entity.ServiceType)(body.ServiceType))
	patchedTender, err = s.tenders.UpdateTender(rctx, patchedTender)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to update tender: %v", err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, api.Tender{
		CreatedAt:      patchedTender.CreatedAt.Format(time.RFC3339),
		Description:    patchedTender.Description,
		Id:             patchedTender.ID.String(),
		Name:           patchedTender.Name,
		OrganizationId: patchedTender.ID.String(),
		ServiceType:    api.TenderServiceType(patchedTender.ServiceType),
		Status:         api.TenderStatus(patchedTender.Status),
		Version:        api.TenderVersion(patchedTender.Version),
	})
}

func (s *Server) RollbackTender(ctx echo.Context, tenderId api.TenderId, version int32, params api.RollbackTenderParams) error {
	rctx := ctx.Request().Context()

	if params.Username == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("add username"),
		})
	}
	if version < 1 {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("invalid version. Version should be >= 1"),
		})
	}

	tenderCreator, err := s.employees.GetByUserName(rctx, params.Username)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get employee by username: %v", err.Error()),
		})
	}

	tenderUUID, err := uuid.Parse(tenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse tenderID: %v", err.Error()),
		})
	}

	tender, err := s.tenders.GetTenderByID(rctx, tenderUUID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender, please, check ID: %v", err.Error()),
		})
	}

	isResponsible, err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, tenderCreator.ID, tender.OrganizationID)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not an organization's responsible: %v", err.Error()),
		})
	}
	if !isResponsible {
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not an organization's responsible"),
		})
	}

	tenderToRollback, err := s.tenders.GetTenderByIDAndVersion(rctx, tenderUUID, int(version))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender's version, please, parametrs: %v", err.Error()),
		})
	}

	tenderToRollback = tender.Rollback(tenderToRollback)

	updatedTender, err := s.tenders.UpdateTender(rctx, tenderToRollback)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to update tender: %v", err.Error()),
		})
	}

	return ctx.JSON(http.StatusOK, api.Tender{
		CreatedAt:      updatedTender.CreatedAt.Format(time.RFC3339),
		Description:    updatedTender.Description,
		Id:             updatedTender.ID.String(),
		Name:           updatedTender.Name,
		OrganizationId: updatedTender.ID.String(),
		ServiceType:    api.TenderServiceType(updatedTender.ServiceType),
		Status:         api.TenderStatus(updatedTender.Status),
		Version:        api.TenderVersion(updatedTender.Version),
	})
}

func (s *Server) GetTenderStatus(ctx echo.Context, tenderId api.TenderId, params api.GetTenderStatusParams) error {
	rctx := ctx.Request().Context()

	var body api.GetTenderStatusParams
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to bind body: %v", err.Error()),
		})
	}

	// поиск пользователя по имени - если нет, то 401
	tenderCreator, err := s.employees.GetByUserName(rctx, *params.Username)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get employee by username: %v", err.Error()),
		})
	}

	tenderIDParsed, err := uuid.Parse(tenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse tenderID: %v", err.Error()),
		})
	}

	tender, err := s.tenders.GetTenderByID(rctx, tenderIDParsed)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender's status by tender ID: %v", err.Error()),
		})
	}

	if tender.Status != entity.TenderStatusPublished {
		// если ID пользователя - не равно ID автора тендера и не в списке ответственных за организацию, то 403
		isResponsible, err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, tenderCreator.ID, tender.OrganizationID)
		if err != nil {
			return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
				Reason: fmt.Sprintf("user is not an organization's responsible: %v", err.Error()),
			})
		}
		if !isResponsible {
			return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
				Reason: fmt.Sprintf("user is not an organization's responsible"),
			})
		}
	}

	return ctx.JSON(http.StatusOK, tender.Status)
}

func (s *Server) UpdateTenderStatus(ctx echo.Context, tenderId api.TenderId, params api.UpdateTenderStatusParams) error {

	// TODO проверить, что статус есть в енаме видов статусов
	rctx := ctx.Request().Context()

	if params.Username == "" || params.Status == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("add status or username"),
		})
	}

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

	tenderUUID, err := uuid.Parse(tenderId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to parse tenderID: %v", err.Error()),
		})
	}

	tender, err := s.tenders.GetTenderByID(rctx, tenderUUID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get tender, please, check ID: %v", err.Error()),
		})
	}

	isResponsible, err := s.organizationResponsibles.IsUserOrganizationResponsible(rctx, employee.ID, tender.OrganizationID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("check if responsible: %v", err.Error()),
		})
	}
	if !isResponsible {
		return ctx.JSON(http.StatusForbidden, api.ErrorResponse{
			Reason: fmt.Sprintf("user is not an organization's responsible"),
		})
	}

	err = s.tenders.UpdateTenderStatus(rctx, tender.ID, entity.TenderStatus(params.Status))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to update tender's status: %v", err.Error()),
		})
	}

	updatedTender, err := s.tenders.GetTenderByID(rctx, tenderUUID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Reason: fmt.Sprintf("failed to get updated tender: %v", err.Error()),
		})
	}
	return ctx.JSON(http.StatusOK, api.Tender{
		CreatedAt:      updatedTender.CreatedAt.Format(time.RFC3339),
		Description:    updatedTender.Description,
		Id:             updatedTender.ID.String(),
		Name:           updatedTender.Name,
		OrganizationId: updatedTender.ID.String(),
		ServiceType:    api.TenderServiceType(updatedTender.ServiceType),
		Status:         api.TenderStatus(updatedTender.Status),
		Version:        api.TenderVersion(updatedTender.Version),
	})
}
