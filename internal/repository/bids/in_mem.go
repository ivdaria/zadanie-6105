package bids

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"sync"
)

type InMemoryRepo struct {
	items map[uuid.UUID]entity.Bid
	mu    sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		items: make(map[uuid.UUID]entity.Bid),
		mu:    sync.RWMutex{},
	}
}

func (r *InMemoryRepo) CreateBid(bid *entity.Bid) uuid.UUID {
	r.mu.Lock()
	defer r.mu.Unlock()

	toSave := *bid
	toSave.ID = uuid.New()
	r.items[toSave.ID] = toSave

	return toSave.ID
}

func (r *InMemoryRepo) GetAllBids() []*entity.Bid {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*entity.Bid, 0, len(r.items))
	for _, v := range r.items {
		result = append(result, &v)
	}

	return result
}

func (r *InMemoryRepo) GetBidsByCreatorID(id uuid.UUID) []*entity.Bid {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Bid
	for _, v := range r.items {
		if v.CreatorID == id {
			result = append(result, &v)
		}
	}

	return result
}

func (r *InMemoryRepo) GetBidsByTenderID(id uuid.UUID) []*entity.Bid {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Bid
	for _, v := range r.items {
		if v.TenderID == id {
			result = append(result, &v)
		}
	}

	return result
}

func (r *InMemoryRepo) GetByID(id uuid.UUID) (*entity.Bid, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bid, ok := r.items[id]
	return &bid, ok
}

func (r *InMemoryRepo) UpdateBid(bid *entity.Bid) {
	r.mu.Lock()
	defer r.mu.Unlock()

	toSave := *bid
	r.items[bid.ID] = toSave
}
