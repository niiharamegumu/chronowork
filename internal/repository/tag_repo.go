package repository

import (
	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/models"
	"gorm.io/gorm"
)

// GormTagRepository is a GORM implementation of TagRepository.
type GormTagRepository struct {
	db *gorm.DB
}

// NewGormTagRepository creates a new GormTagRepository.
func NewGormTagRepository(db *gorm.DB) *GormTagRepository {
	return &GormTagRepository{db: db}
}

// Create creates a new Tag.
func (r *GormTagRepository) Create(name string) (*domain.Tag, error) {
	tag := models.Tag{Name: name}
	if err := r.db.Create(&tag).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&tag), nil
}

// FindByID finds a Tag by its ID.
func (r *GormTagRepository) FindByID(id uint) (*domain.Tag, error) {
	var tag models.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&tag), nil
}

// FindAll finds all Tags.
func (r *GormTagRepository) FindAll() ([]domain.Tag, error) {
	var tags []models.Tag
	if err := r.db.Order("id desc").Find(&tags).Error; err != nil {
		return nil, err
	}
	return r.toDomainSlice(tags), nil
}

// FindByNames finds Tags by their names.
func (r *GormTagRepository) FindByNames(names []string) ([]domain.Tag, error) {
	var tags []models.Tag
	if err := r.db.Where("name IN ?", names).Find(&tags).Error; err != nil {
		return nil, err
	}
	return r.toDomainSlice(tags), nil
}

// GetAllNames returns all Tag names.
func (r *GormTagRepository) GetAllNames() []string {
	var tags []models.Tag
	if err := r.db.Find(&tags).Error; err != nil {
		return []string{}
	}
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames
}

// Update updates a Tag's name.
func (r *GormTagRepository) Update(id uint, name string) error {
	return r.db.Model(&models.Tag{}).Where("id = ?", id).Update("name", name).Error
}

// Delete permanently deletes a Tag.
func (r *GormTagRepository) Delete(id uint) error {
	return r.db.Unscoped().Delete(&models.Tag{}, id).Error
}

// toDomain converts a GORM model to a domain entity.
func (r *GormTagRepository) toDomain(m *models.Tag) *domain.Tag {
	return &domain.Tag{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// toDomainSlice converts a slice of GORM models to domain entities.
func (r *GormTagRepository) toDomainSlice(ms []models.Tag) []domain.Tag {
	ds := make([]domain.Tag, len(ms))
	for i, m := range ms {
		ds[i] = *r.toDomain(&m)
	}
	return ds
}
