package gateway

import (
	"context"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
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
	GetBidsByTenderID(ctx context.Context, tenderID uuid.UUID, pagination entity.Pagination) ([]*entity.Bid, error)
	GetBidByID(ctx context.Context, id uuid.UUID) (*entity.Bid, error)
	UpdateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error)
	GetBidByIDAndVersion(ctx context.Context, id uuid.UUID, version int) (*entity.Bid, error)
	GetBidsByUsername(ctx context.Context, username string, pagination entity.Pagination) ([]*entity.Bid, error)
	UpdateBidDecision(ctx context.Context, id uuid.UUID, bidDecision entity.BidDecision) error
	UpdateBidStatus(ctx context.Context, id uuid.UUID, newStatus entity.BidStatus) error
}

type organizationResponsibleRepo interface {
	IsUserOrganizationResponsible(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) error
	IsUserResponsible(ctx context.Context, userID uuid.UUID) error
	GetOrganizationResponsibleByUserID(ctx context.Context, userID uuid.UUID) (*entity.OrganizationResponsible, error)
}

type userCanEditBidChecker interface {
	IsUserCanEditBid(ctx context.Context, bid *entity.Bid, user *entity.Employee) (bool, error)
}

type Server struct {
	tenders                  tenderRepo
	employees                employeeRepo
	organizations            organizationRepo
	bids                     bidRepo
	organizationResponsibles organizationResponsibleRepo
	userCanEditBidChecker    userCanEditBidChecker
}

func NewServer(
	tenders tenderRepo,
	employees employeeRepo,
	organizations organizationRepo,
	bids bidRepo,
	organizationResponsibles organizationResponsibleRepo,
	userCanEditBidChecker userCanEditBidChecker,
) *Server {
	return &Server{
		tenders:                  tenders,
		employees:                employees,
		organizations:            organizations,
		bids:                     bids,
		organizationResponsibles: organizationResponsibles,
		userCanEditBidChecker:    userCanEditBidChecker,
	}
}
