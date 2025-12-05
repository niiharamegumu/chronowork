package mock

import (
	"errors"
	"sync"
	"time"

	"github.com/niiharamegumu/chronowork/internal/domain"
)

// ProjectTypeRepository is an in-memory mock of repository.ProjectTypeRepository.
type ProjectTypeRepository struct {
	mu      sync.RWMutex
	data    map[uint]*domain.ProjectType
	nextID  uint
	tagRepo *TagRepository
}

// NewProjectTypeRepository creates a new mock ProjectTypeRepository.
func NewProjectTypeRepository(tagRepo *TagRepository) *ProjectTypeRepository {
	return &ProjectTypeRepository{
		data:    make(map[uint]*domain.ProjectType),
		nextID:  1,
		tagRepo: tagRepo,
	}
}

// Create creates a new ProjectType with optional tags.
func (r *ProjectTypeRepository) Create(name string, tagIDs []uint) (*domain.ProjectType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	pt := &domain.ProjectType{
		ID:        r.nextID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Add tags if tagRepo is available
	if r.tagRepo != nil && len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			tag, err := r.tagRepo.FindByID(tagID)
			if err == nil {
				pt.Tags = append(pt.Tags, *tag)
			}
		}
	}

	r.data[r.nextID] = pt
	r.nextID++
	return pt, nil
}

// FindByID finds a ProjectType by its ID with tags preloaded.
func (r *ProjectTypeRepository) FindByID(id uint) (*domain.ProjectType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pt, ok := r.data[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return pt, nil
}

// FindByName finds a ProjectType by its name with tags preloaded.
func (r *ProjectTypeRepository) FindByName(name string) (*domain.ProjectType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, pt := range r.data {
		if pt.Name == name {
			return pt, nil
		}
	}
	return &domain.ProjectType{}, nil
}

// FindAllWithTags finds all ProjectTypes with tags preloaded.
func (r *ProjectTypeRepository) FindAllWithTags() ([]domain.ProjectType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.ProjectType
	for _, pt := range r.data {
		result = append(result, *pt)
	}
	return result, nil
}

// GetAllNames returns all ProjectType names.
func (r *ProjectTypeRepository) GetAllNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var names []string
	for _, pt := range r.data {
		names = append(names, pt.Name)
	}
	return names
}

// Update updates a ProjectType's name and tags.
func (r *ProjectTypeRepository) Update(id uint, name string, tagIDs []uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	pt, ok := r.data[id]
	if !ok {
		return errors.New("record not found")
	}

	pt.Name = name
	pt.Tags = nil
	pt.UpdatedAt = time.Now()

	if r.tagRepo != nil && len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			tag, err := r.tagRepo.FindByID(tagID)
			if err == nil {
				pt.Tags = append(pt.Tags, *tag)
			}
		}
	}

	return nil
}

// Delete permanently deletes a ProjectType.
func (r *ProjectTypeRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return errors.New("record not found")
	}
	delete(r.data, id)
	return nil
}
