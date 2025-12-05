// Package usecase contains business logic implementations.
package usecase

import (
	"errors"
	"time"

	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/internal/repository"
)

// ChronoWorkUseCase handles business logic for ChronoWork operations.
type ChronoWorkUseCase struct {
	repo repository.ChronoWorkRepository
}

// NewChronoWorkUseCase creates a new ChronoWorkUseCase.
func NewChronoWorkUseCase(repo repository.ChronoWorkRepository) *ChronoWorkUseCase {
	return &ChronoWorkUseCase{repo: repo}
}

// Create creates a new ChronoWork entry.
func (uc *ChronoWorkUseCase) Create(title string, projectTypeID, tagID uint) (*domain.ChronoWork, error) {
	// Check for duplicate title today
	existing, err := uc.repo.FindByTitleToday(title)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("work with this title already exists today")
	}

	return uc.repo.Create(title, projectTypeID, tagID)
}

// FindByID finds a ChronoWork by its ID.
func (uc *ChronoWorkUseCase) FindByID(id uint) (*domain.ChronoWork, error) {
	return uc.repo.FindByID(id)
}

// FindInRange finds ChronoWorks within a time range.
func (uc *ChronoWorkUseCase) FindInRange(startTime, endTime time.Time) ([]domain.ChronoWork, error) {
	return uc.repo.FindInRange(startTime, endTime)
}

// FindTracking finds all currently tracking ChronoWorks.
func (uc *ChronoWorkUseCase) FindTracking() ([]domain.ChronoWork, error) {
	return uc.repo.FindTracking()
}

// FindByProjectTypeID finds ChronoWorks by project type ID.
func (uc *ChronoWorkUseCase) FindByProjectTypeID(projectTypeID uint) ([]domain.ChronoWork, error) {
	return uc.repo.FindByProjectTypeID(projectTypeID)
}

// GetAll returns all ChronoWorks with optional ordering and limit.
func (uc *ChronoWorkUseCase) GetAll(orderField string, limit int) ([]domain.ChronoWork, error) {
	return uc.repo.GetAll(orderField, limit)
}

// Update updates a ChronoWork's title, projectTypeID, and tagID.
func (uc *ChronoWorkUseCase) Update(id uint, title string, projectTypeID, tagID uint) error {
	return uc.repo.Update(id, title, projectTypeID, tagID)
}

// UpdateTotalSeconds updates the total seconds of a ChronoWork.
func (uc *ChronoWorkUseCase) UpdateTotalSeconds(id uint, totalSeconds int) error {
	return uc.repo.UpdateTotalSeconds(id, totalSeconds)
}

// UpdateConfirmed updates the confirmed status of a ChronoWork.
func (uc *ChronoWorkUseCase) UpdateConfirmed(id uint, confirmed bool) error {
	return uc.repo.UpdateConfirmed(id, confirmed)
}

// StartTracking starts tracking a ChronoWork.
func (uc *ChronoWorkUseCase) StartTracking(id uint) error {
	return uc.repo.StartTracking(id)
}

// StopTracking stops tracking a ChronoWork and calculates total time.
func (uc *ChronoWorkUseCase) StopTracking(id uint) error {
	return uc.repo.StopTracking(id)
}

// Delete permanently deletes a ChronoWork.
func (uc *ChronoWorkUseCase) Delete(id uint) error {
	return uc.repo.Delete(id)
}
