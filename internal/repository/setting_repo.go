package repository

import (
	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/models"
	"gorm.io/gorm"
)

// GormSettingRepository is a GORM implementation of SettingRepository.
type GormSettingRepository struct {
	db *gorm.DB
}

// NewGormSettingRepository creates a new GormSettingRepository.
func NewGormSettingRepository(db *gorm.DB) *GormSettingRepository {
	return &GormSettingRepository{db: db}
}

// Get retrieves the current setting or creates a default one.
func (r *GormSettingRepository) Get() (*domain.Setting, error) {
	var setting models.Setting
	if err := r.db.FirstOrCreate(&setting).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&setting), nil
}

// Update updates the setting.
func (r *GormSettingRepository) Update(setting *domain.Setting) error {
	return r.db.Model(&models.Setting{}).Where("id = ?", setting.ID).
		Select("relative_date", "person_day", "display_as_person_day", "download_path").
		Updates(map[string]interface{}{
			"relative_date":         setting.RelativeDate,
			"person_day":            setting.PersonDay,
			"display_as_person_day": setting.DisplayAsPersonDay,
			"download_path":         setting.DownloadPath,
		}).Error
}

// toDomain converts a GORM model to a domain entity.
func (r *GormSettingRepository) toDomain(m *models.Setting) *domain.Setting {
	return &domain.Setting{
		ID:                 m.ID,
		RelativeDate:       m.RelativeDate,
		PersonDay:          m.PersonDay,
		DisplayAsPersonDay: m.DisplayAsPersonDay,
		DownloadPath:       m.DownloadPath,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}
