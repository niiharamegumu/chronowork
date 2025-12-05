package usecase

import (
	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/internal/repository"
)

// SettingUseCase handles business logic for Setting operations.
type SettingUseCase struct {
	repo repository.SettingRepository
}

// NewSettingUseCase creates a new SettingUseCase.
func NewSettingUseCase(repo repository.SettingRepository) *SettingUseCase {
	return &SettingUseCase{repo: repo}
}

// Get retrieves the current setting or creates a default one.
func (uc *SettingUseCase) Get() (*domain.Setting, error) {
	return uc.repo.Get()
}

// Update updates the setting.
func (uc *SettingUseCase) Update(setting *domain.Setting) error {
	return uc.repo.Update(setting)
}
