package mock

import (
	"errors"
	"sync"
	"time"

	"github.com/niiharamegumu/chronowork/internal/domain"
)

// TagRepository is an in-memory mock of repository.TagRepository.
type TagRepository struct {
	mu     sync.RWMutex
	data   map[uint]*domain.Tag
	nextID uint
}

// NewTagRepository creates a new mock TagRepository.
func NewTagRepository() *TagRepository {
	return &TagRepository{
		data:   make(map[uint]*domain.Tag),
		nextID: 1,
	}
}

// Create creates a new Tag.
func (r *TagRepository) Create(name string) (*domain.Tag, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	tag := &domain.Tag{
		ID:        r.nextID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.data[r.nextID] = tag
	r.nextID++
	return tag, nil
}

// FindByID finds a Tag by its ID.
func (r *TagRepository) FindByID(id uint) (*domain.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tag, ok := r.data[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return tag, nil
}

// FindAll finds all Tags.
func (r *TagRepository) FindAll() ([]domain.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Tag
	for _, tag := range r.data {
		result = append(result, *tag)
	}
	return result, nil
}

// FindByNames finds Tags by their names.
func (r *TagRepository) FindByNames(names []string) ([]domain.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}

	var result []domain.Tag
	for _, tag := range r.data {
		if nameSet[tag.Name] {
			result = append(result, *tag)
		}
	}
	return result, nil
}

// GetAllNames returns all Tag names.
func (r *TagRepository) GetAllNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var names []string
	for _, tag := range r.data {
		names = append(names, tag.Name)
	}
	return names
}

// Update updates a Tag's name.
func (r *TagRepository) Update(id uint, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tag, ok := r.data[id]
	if !ok {
		return errors.New("record not found")
	}
	tag.Name = name
	tag.UpdatedAt = time.Now()
	return nil
}

// Delete permanently deletes a Tag.
func (r *TagRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return errors.New("record not found")
	}
	delete(r.data, id)
	return nil
}
