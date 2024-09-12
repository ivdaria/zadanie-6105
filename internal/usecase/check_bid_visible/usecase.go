package check_bid_visible

import (
	"context"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
)

type UseCase struct {
}

func (uc *UseCase) IsBidVisibleToUser(ctx context.Context, bid *entity.Bid, user *entity.Employee) (bool, error) {

	return true, nil
}
