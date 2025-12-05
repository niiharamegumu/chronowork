package usecase

import (
	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/internal/repository"
)

// ProjectTypeUseCase handles business logic for ProjectType operations.
type ProjectTypeUseCase struct {
	repo repository.ProjectTypeRepository
}

// NewProjectTypeUseCase creates a new ProjectTypeUseCase.
func NewProjectTypeUseCase(repo repository.ProjectTypeRepository) *ProjectTypeUseCase {
	return &ProjectTypeUseCase{repo: repo}
}

// Create creates a new ProjectType with optional tags.
func (uc *ProjectTypeUseCase) Create(name string, tagIDs []uint) (*domain.ProjectType, error) {
	return uc.repo.Create(name, tagIDs)
}

// FindByID finds a ProjectType by its ID with tags preloaded.
func (uc *ProjectTypeUseCase) FindByID(id uint) (*domain.ProjectType, error) {
	return uc.repo.FindByID(id)
}

// FindByName finds a ProjectType by its name with tags preloaded.
func (uc *ProjectTypeUseCase) FindByName(name string) (*domain.ProjectType, error) {
	return uc.repo.FindByName(name)
}

// FindAllWithTags finds all ProjectTypes with tags preloaded.
func (uc *ProjectTypeUseCase) FindAllWithTags() ([]domain.ProjectType, error) {
	return uc.repo.FindAllWithTags()
}

// GetAllNames returns all ProjectType names.
func (uc *ProjectTypeUseCase) GetAllNames() []string {
	return uc.repo.GetAllNames()
}

// Update updates a ProjectType's name and tags.
func (uc *ProjectTypeUseCase) Update(id uint, name string, tagIDs []uint) error {
	return uc.repo.Update(id, name, tagIDs)
}

// Delete permanently deletes a ProjectType.
func (uc *ProjectTypeUseCase) Delete(id uint) error {
	return uc.repo.Delete(id)
}
