package mock

import (
	"sync"
	"time"

	"github.com/niiharamegumu/chronowork/internal/domain"
)

// SettingRepository is an in-memory mock of repository.SettingRepository.
type SettingRepository struct {
	mu      sync.RWMutex
	setting *domain.Setting
}

// NewSettingRepository creates a new mock SettingRepository.
func NewSettingRepository() *SettingRepository {
	return &SettingRepository{}
}

// Get retrieves the current setting or creates a default one.
func (r *SettingRepository) Get() (*domain.Setting, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.setting == nil {
		now := time.Now()
		r.setting = &domain.Setting{
			ID:                 1,
			RelativeDate:       0,
			PersonDay:          8,
			DisplayAsPersonDay: true,
			DownloadPath:       "./",
			CreatedAt:          now,
			UpdatedAt:          now,
		}
	}
	return r.setting, nil
}

// Update updates the setting.
func (r *SettingRepository) Update(setting *domain.Setting) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.setting == nil {
		now := time.Now()
		r.setting = &domain.Setting{
			ID:        1,
			CreatedAt: now,
		}
	}
	r.setting.RelativeDate = setting.RelativeDate
	r.setting.PersonDay = setting.PersonDay
	r.setting.DisplayAsPersonDay = setting.DisplayAsPersonDay
	r.setting.DownloadPath = setting.DownloadPath
	r.setting.UpdatedAt = time.Now()
	return nil
}
