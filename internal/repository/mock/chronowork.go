// Package mock provides in-memory mock implementations of repository interfaces for testing.
package mock

import (
	"errors"
	"sync"
	"time"

	"github.com/niiharamegumu/chronowork/internal/domain"
)

// ChronoWorkRepository is an in-memory mock of repository.ChronoWorkRepository.
type ChronoWorkRepository struct {
	mu          sync.RWMutex
	data        map[uint]*domain.ChronoWork
	nextID      uint
	findByIDErr error
}

// NewChronoWorkRepository creates a new mock ChronoWorkRepository.
func NewChronoWorkRepository() *ChronoWorkRepository {
	return &ChronoWorkRepository{
		data:   make(map[uint]*domain.ChronoWork),
		nextID: 1,
	}
}

// SetFindByIDError sets an error to be returned by FindByID (for testing error cases).
func (r *ChronoWorkRepository) SetFindByIDError(err error) {
	r.findByIDErr = err
}

// Create creates a new ChronoWork entry.
func (r *ChronoWorkRepository) Create(title string, projectTypeID, tagID uint) (*domain.ChronoWork, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cw := &domain.ChronoWork{
		ID:            r.nextID,
		Title:         title,
		ProjectTypeID: projectTypeID,
		TagID:         tagID,
		StartTime:     time.Time{},
		EndTime:       time.Time{},
		IsTracking:    false,
		TotalSeconds:  0,
		Confirmed:     false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	r.data[r.nextID] = cw
	r.nextID++
	return cw, nil
}

// FindByID finds a ChronoWork by its ID.
func (r *ChronoWorkRepository) FindByID(id uint) (*domain.ChronoWork, error) {
	if r.findByIDErr != nil {
		return nil, r.findByIDErr
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	cw, ok := r.data[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return cw, nil
}

// FindInRange finds ChronoWorks within a time range.
func (r *ChronoWorkRepository) FindInRange(startTime, endTime time.Time) ([]domain.ChronoWork, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.ChronoWork
	for _, cw := range r.data {
		if cw.CreatedAt.After(startTime) && cw.CreatedAt.Before(endTime) {
			result = append(result, *cw)
		}
	}
	return result, nil
}

// FindTracking finds all currently tracking ChronoWorks.
func (r *ChronoWorkRepository) FindTracking() ([]domain.ChronoWork, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.ChronoWork
	for _, cw := range r.data {
		if cw.IsTracking {
			result = append(result, *cw)
		}
	}
	return result, nil
}

// FindByTitleToday finds a ChronoWork by title created today.
func (r *ChronoWorkRepository) FindByTitleToday(title string) (*domain.ChronoWork, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, time.Local)

	for _, cw := range r.data {
		if cw.Title == title && cw.CreatedAt.After(startOfDay) && cw.CreatedAt.Before(endOfDay) {
			return cw, nil
		}
	}
	return nil, nil
}

// FindByProjectTypeID finds ChronoWorks by project type ID.
func (r *ChronoWorkRepository) FindByProjectTypeID(projectTypeID uint) ([]domain.ChronoWork, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.ChronoWork
	for _, cw := range r.data {
		if cw.ProjectTypeID == projectTypeID {
			result = append(result, *cw)
		}
	}
	return result, nil
}

// GetAll returns all ChronoWorks with optional ordering and limit.
func (r *ChronoWorkRepository) GetAll(orderField string, limit int) ([]domain.ChronoWork, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.ChronoWork
	count := 0
	for _, cw := range r.data {
		if limit > 0 && count >= limit {
			break
		}
		result = append(result, *cw)
		count++
	}
	return result, nil
}

// Update updates a ChronoWork's title, projectTypeID, and tagID.
func (r *ChronoWorkRepository) Update(id uint, title string, projectTypeID, tagID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cw, ok := r.data[id]
	if !ok {
		return errors.New("record not found")
	}
	cw.Title = title
	cw.ProjectTypeID = projectTypeID
	cw.TagID = tagID
	cw.UpdatedAt = time.Now()
	return nil
}

// UpdateTotalSeconds updates the total seconds of a ChronoWork.
func (r *ChronoWorkRepository) UpdateTotalSeconds(id uint, totalSeconds int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cw, ok := r.data[id]
	if !ok {
		return errors.New("record not found")
	}
	cw.TotalSeconds = totalSeconds
	cw.UpdatedAt = time.Now()
	return nil
}

// UpdateConfirmed updates the confirmed status of a ChronoWork.
func (r *ChronoWorkRepository) UpdateConfirmed(id uint, confirmed bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cw, ok := r.data[id]
	if !ok {
		return errors.New("record not found")
	}
	cw.Confirmed = confirmed
	cw.UpdatedAt = time.Now()
	return nil
}

// StartTracking starts tracking a ChronoWork.
func (r *ChronoWorkRepository) StartTracking(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cw, ok := r.data[id]
	if !ok {
		return errors.New("record not found")
	}
	cw.StartTime = time.Now()
	cw.EndTime = time.Time{}
	cw.IsTracking = true
	cw.UpdatedAt = time.Now()
	return nil
}

// StopTracking stops tracking a ChronoWork and calculates total time.
func (r *ChronoWorkRepository) StopTracking(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cw, ok := r.data[id]
	if !ok {
		return errors.New("record not found")
	}
	cw.EndTime = time.Now()
	cw.IsTracking = false
	elapsed := int(cw.EndTime.Sub(cw.StartTime).Seconds())
	cw.TotalSeconds += elapsed
	cw.UpdatedAt = time.Now()
	return nil
}

// Delete permanently deletes a ChronoWork.
func (r *ChronoWorkRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return errors.New("record not found")
	}
	delete(r.data, id)
	return nil
}
