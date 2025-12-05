package repository

import (
	"math"
	"time"

	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/models"
	"gorm.io/gorm"
)

// GormChronoWorkRepository is a GORM implementation of ChronoWorkRepository.
type GormChronoWorkRepository struct {
	db *gorm.DB
}

// NewGormChronoWorkRepository creates a new GormChronoWorkRepository.
func NewGormChronoWorkRepository(db *gorm.DB) *GormChronoWorkRepository {
	return &GormChronoWorkRepository{db: db}
}

// Create creates a new ChronoWork entry.
func (r *GormChronoWorkRepository) Create(title string, projectTypeID, tagID uint) (*domain.ChronoWork, error) {
	chronoWork := models.ChronoWork{
		Title:         title,
		ProjectTypeID: projectTypeID,
		TagID:         tagID,
		StartTime:     time.Time{},
		EndTime:       time.Time{},
		IsTracking:    false,
		TotalSeconds:  0,
		Confirmed:     false,
	}
	if err := r.db.Create(&chronoWork).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&chronoWork), nil
}

// FindByID finds a ChronoWork by its ID.
func (r *GormChronoWorkRepository) FindByID(id uint) (*domain.ChronoWork, error) {
	var chronoWork models.ChronoWork
	if err := r.db.Preload("ProjectType").Preload("Tag").First(&chronoWork, id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&chronoWork), nil
}

// FindInRange finds ChronoWorks within a time range.
func (r *GormChronoWorkRepository) FindInRange(startTime, endTime time.Time) ([]domain.ChronoWork, error) {
	var chronoWorks []models.ChronoWork
	err := r.db.
		Preload("ProjectType").
		Preload("Tag").
		Order("created_at desc").
		Order("id desc").
		Find(&chronoWorks, "created_at >= ? AND created_at <= ?", startTime, endTime).Error
	if err != nil {
		return nil, err
	}
	return r.toDomainSlice(chronoWorks), nil
}

// FindTracking finds all currently tracking ChronoWorks.
func (r *GormChronoWorkRepository) FindTracking() ([]domain.ChronoWork, error) {
	var chronoWorks []models.ChronoWork
	err := r.db.
		Preload("ProjectType").
		Preload("Tag").
		Find(&chronoWorks, "is_tracking = ?", true).Error
	if err != nil {
		return nil, err
	}
	return r.toDomainSlice(chronoWorks), nil
}

// FindByProjectTypeID finds ChronoWorks by project type ID.
func (r *GormChronoWorkRepository) FindByProjectTypeID(projectTypeID uint) ([]domain.ChronoWork, error) {
	var chronoWorks []models.ChronoWork
	err := r.db.
		Preload("ProjectType").
		Find(&chronoWorks, "project_type_id = ?", projectTypeID).Error
	if err != nil {
		return nil, err
	}
	return r.toDomainSlice(chronoWorks), nil
}

// GetAll returns all ChronoWorks with optional ordering and limit.
func (r *GormChronoWorkRepository) GetAll(orderField string, limit int) ([]domain.ChronoWork, error) {
	var chronoWorks []models.ChronoWork
	query := r.db.Preload("ProjectType").Preload("Tag")
	if orderField != "" {
		query = query.Order(orderField)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&chronoWorks).Error; err != nil {
		return nil, err
	}
	return r.toDomainSlice(chronoWorks), nil
}

// Update updates a ChronoWork's title, projectTypeID, and tagID.
func (r *GormChronoWorkRepository) Update(id uint, title string, projectTypeID, tagID uint) error {
	return r.db.Model(&models.ChronoWork{}).Where("id = ?", id).
		Select("title", "project_type_id", "tag_id").
		Updates(map[string]interface{}{
			"title":           title,
			"project_type_id": projectTypeID,
			"tag_id":          tagID,
		}).Error
}

// UpdateTotalSeconds updates the total seconds of a ChronoWork.
func (r *GormChronoWorkRepository) UpdateTotalSeconds(id uint, totalSeconds int) error {
	return r.db.Model(&models.ChronoWork{}).Where("id = ?", id).
		Select("total_seconds").
		Updates(map[string]interface{}{
			"total_seconds": totalSeconds,
		}).Error
}

// UpdateConfirmed updates the confirmed status of a ChronoWork.
func (r *GormChronoWorkRepository) UpdateConfirmed(id uint, confirmed bool) error {
	return r.db.Model(&models.ChronoWork{}).Where("id = ?", id).
		Select("Confirmed").
		Updates(map[string]interface{}{
			"confirmed": confirmed,
		}).Error
}

// StartTracking starts tracking a ChronoWork.
func (r *GormChronoWorkRepository) StartTracking(id uint) error {
	return r.db.Model(&models.ChronoWork{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"start_time":  time.Now(),
			"end_time":    time.Time{},
			"is_tracking": true,
		}).Error
}

// StopTracking stops tracking a ChronoWork and calculates total time.
func (r *GormChronoWorkRepository) StopTracking(id uint) error {
	var chronoWork models.ChronoWork
	if err := r.db.First(&chronoWork, id).Error; err != nil {
		return err
	}

	endTime := time.Now()
	elapsed := int(math.Floor(endTime.Sub(chronoWork.StartTime).Seconds()))

	return r.db.Model(&models.ChronoWork{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"end_time":      endTime,
			"is_tracking":   false,
			"total_seconds": gorm.Expr("total_seconds + ?", elapsed),
		}).Error
}

// Delete permanently deletes a ChronoWork.
func (r *GormChronoWorkRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&models.ChronoWork{}, id).Error
}

// toDomain converts a GORM model to a domain entity.
func (r *GormChronoWorkRepository) toDomain(m *models.ChronoWork) *domain.ChronoWork {
	d := &domain.ChronoWork{
		ID:            m.ID,
		Title:         m.Title,
		ProjectTypeID: m.ProjectTypeID,
		TagID:         m.TagID,
		StartTime:     m.StartTime,
		EndTime:       m.EndTime,
		IsTracking:    m.IsTracking,
		TotalSeconds:  m.TotalSeconds,
		Confirmed:     m.Confirmed,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
	if m.ProjectType.ID != 0 {
		d.ProjectType = &domain.ProjectType{
			ID:   m.ProjectType.ID,
			Name: m.ProjectType.Name,
		}
	}
	if m.Tag.ID != 0 {
		d.Tag = &domain.Tag{
			ID:   m.Tag.ID,
			Name: m.Tag.Name,
		}
	}
	return d
}

// toDomainSlice converts a slice of GORM models to domain entities.
func (r *GormChronoWorkRepository) toDomainSlice(ms []models.ChronoWork) []domain.ChronoWork {
	ds := make([]domain.ChronoWork, len(ms))
	for i, m := range ms {
		ds[i] = *r.toDomain(&m)
	}
	return ds
}
