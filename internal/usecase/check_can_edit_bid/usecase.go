package check_can_edit_bid

import (
	"context"
	"errors"
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type organizationResponsibleRepo interface {
	IsUserOrganizationResponsible(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) error
	GetOrganizationResponsibleByUserID(ctx context.Context, userID uuid.UUID) (*entity.OrganizationResponsible, error)
}

type UseCase struct {
	organizationResponsibles organizationResponsibleRepo
}

func NewUseCase(organizationResponsibles organizationResponsibleRepo) *UseCase {
	return &UseCase{organizationResponsibles: organizationResponsibles}
}

func (uc *UseCase) IsUserCanEditBid(ctx context.Context, bid *entity.Bid, user *entity.Employee) (bool, error) {
	// если автор предложения - пользователь (user), то откатывать версию предложения может только он
	// если автор предложения - организация (organization), то только ответственные пользователи организации
	if bid.AuthorType == entity.BidAuthorTypeUser {
		if bid.CreatorID != user.ID {
			return false, nil
		}
	} else {
		// проверяем, является ли пользователь ответственным
		oldBitOrganization, err := uc.organizationResponsibles.GetOrganizationResponsibleByUserID(ctx, bid.CreatorID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, nil
			}
			return false, fmt.Errorf("GetOrganizationResponsibleByUserID: %v", err)
		}

		if err := uc.organizationResponsibles.IsUserOrganizationResponsible(ctx, user.ID, oldBitOrganization.OrganizationID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, nil
			}
			return false, fmt.Errorf("IsUserOrganizationResponsible: %v", err)
		}
	}

	return true, nil
}
