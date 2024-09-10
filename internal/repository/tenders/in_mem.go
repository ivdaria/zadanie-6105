package tenders

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725731644-team-78845/zadanie-6105/internal/entity"
	"github.com/google/uuid"
	"sync"
	"time"
)

// логика хранения тендеров
type InMemoryRepo struct {
	items map[uuid.UUID]entity.Tender
	mu    sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		items: make(map[uuid.UUID]entity.Tender),
		mu:    sync.RWMutex{},
	}
}

func (r *InMemoryRepo) CreateTender(tender *entity.Tender) *entity.Tender {
	r.mu.Lock()
	defer r.mu.Unlock()

	toSave := *tender
	toSave.ID = uuid.New()
	toSave.CreatedAt = time.Now()
	r.items[toSave.ID] = toSave

	return &toSave
}

func (r *InMemoryRepo) UpdateTender(tender *entity.Tender) {
	r.mu.Lock()
	defer r.mu.Unlock()

	toSave := *tender
	r.items[tender.ID] = toSave
}

func (r *InMemoryRepo) DeleteTender(id uuid.UUID) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; ok {
		delete(r.items, id)
		return true
	}
	return false
}

func (r *InMemoryRepo) GetAllTenders() []*entity.Tender {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*entity.Tender, 0, len(r.items))
	for _, v := range r.items {
		result = append(result, &v)
	}

	return result
}

func (r *InMemoryRepo) GetTendersByUserID(id uuid.UUID) []*entity.Tender {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Tender
	for _, v := range r.items {
		if v.CreatorID == id {
			result = append(result, &v)
		}
	}

	return result
}

//откат версии
