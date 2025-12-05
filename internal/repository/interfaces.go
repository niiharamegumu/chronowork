// Package repository defines repository interfaces and their implementations.
package repository

import (
	"time"

	"github.com/niiharamegumu/chronowork/internal/domain"
)

// ChronoWorkRepository defines operations for ChronoWork persistence.
type ChronoWorkRepository interface {
	// Create creates a new ChronoWork entry.
	Create(title string, projectTypeID, tagID uint) (*domain.ChronoWork, error)
	// FindByID finds a ChronoWork by its ID.
	FindByID(id uint) (*domain.ChronoWork, error)
	// FindInRange finds ChronoWorks within a time range.
	FindInRange(startTime, endTime time.Time) ([]domain.ChronoWork, error)
	// FindTracking finds all currently tracking ChronoWorks.
	FindTracking() ([]domain.ChronoWork, error)
	// FindByProjectTypeID finds ChronoWorks by project type ID.
	FindByProjectTypeID(projectTypeID uint) ([]domain.ChronoWork, error)
	// GetAll returns all ChronoWorks with optional ordering and limit.
	GetAll(orderField string, limit int) ([]domain.ChronoWork, error)
	// Update updates a ChronoWork's title, projectTypeID, and tagID.
	Update(id uint, title string, projectTypeID, tagID uint) error
	// UpdateTotalSeconds updates the total seconds of a ChronoWork.
	UpdateTotalSeconds(id uint, totalSeconds int) error
	// UpdateConfirmed updates the confirmed status of a ChronoWork.
	UpdateConfirmed(id uint, confirmed bool) error
	// StartTracking starts tracking a ChronoWork.
	StartTracking(id uint) error
	// StopTracking stops tracking a ChronoWork and calculates total time.
	StopTracking(id uint) error
	// Delete permanently deletes a ChronoWork.
	Delete(id uint) error
}

// TagRepository defines operations for Tag persistence.
type TagRepository interface {
	// Create creates a new Tag.
	Create(name string) (*domain.Tag, error)
	// FindByID finds a Tag by its ID.
	FindByID(id uint) (*domain.Tag, error)
	// FindAll finds all Tags.
	FindAll() ([]domain.Tag, error)
	// FindByNames finds Tags by their names.
	FindByNames(names []string) ([]domain.Tag, error)
	// GetAllNames returns all Tag names.
	GetAllNames() []string
	// Update updates a Tag's name.
	Update(id uint, name string) error
	// Delete permanently deletes a Tag.
	Delete(id uint) error
}

// ProjectTypeRepository defines operations for ProjectType persistence.
type ProjectTypeRepository interface {
	// Create creates a new ProjectType with optional tags.
	Create(name string, tagIDs []uint) (*domain.ProjectType, error)
	// FindByID finds a ProjectType by its ID with tags preloaded.
	FindByID(id uint) (*domain.ProjectType, error)
	// FindByName finds a ProjectType by its name with tags preloaded.
	FindByName(name string) (*domain.ProjectType, error)
	// FindAllWithTags finds all ProjectTypes with tags preloaded.
	FindAllWithTags() ([]domain.ProjectType, error)
	// GetAllNames returns all ProjectType names.
	GetAllNames() []string
	// Update updates a ProjectType's name and tags.
	Update(id uint, name string, tagIDs []uint) error
	// Delete permanently deletes a ProjectType.
	Delete(id uint) error
}

// SettingRepository defines operations for Setting persistence.
type SettingRepository interface {
	// Get retrieves the current setting or creates a default one.
	Get() (*domain.Setting, error)
	// Update updates the setting.
	Update(setting *domain.Setting) error
}
