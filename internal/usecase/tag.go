package usecase

import (
	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/internal/repository"
)

// TagUseCase handles business logic for Tag operations.
type TagUseCase struct {
	repo repository.TagRepository
}

// NewTagUseCase creates a new TagUseCase.
func NewTagUseCase(repo repository.TagRepository) *TagUseCase {
	return &TagUseCase{repo: repo}
}

// Create creates a new Tag.
func (uc *TagUseCase) Create(name string) (*domain.Tag, error) {
	return uc.repo.Create(name)
}

// FindByID finds a Tag by its ID.
func (uc *TagUseCase) FindByID(id uint) (*domain.Tag, error) {
	return uc.repo.FindByID(id)
}

// FindAll finds all Tags.
func (uc *TagUseCase) FindAll() ([]domain.Tag, error) {
	return uc.repo.FindAll()
}

// FindByNames finds Tags by their names.
func (uc *TagUseCase) FindByNames(names []string) ([]domain.Tag, error) {
	return uc.repo.FindByNames(names)
}

// GetAllNames returns all Tag names.
func (uc *TagUseCase) GetAllNames() []string {
	return uc.repo.GetAllNames()
}

// Update updates a Tag's name.
func (uc *TagUseCase) Update(id uint, name string) error {
	return uc.repo.Update(id, name)
}

// Delete permanently deletes a Tag.
func (uc *TagUseCase) Delete(id uint) error {
	return uc.repo.Delete(id)
}
