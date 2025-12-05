package repository

import (
	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/models"
	"gorm.io/gorm"
)

// GormProjectTypeRepository is a GORM implementation of ProjectTypeRepository.
type GormProjectTypeRepository struct {
	db *gorm.DB
}

// NewGormProjectTypeRepository creates a new GormProjectTypeRepository.
func NewGormProjectTypeRepository(db *gorm.DB) *GormProjectTypeRepository {
	return &GormProjectTypeRepository{db: db}
}

// Create creates a new ProjectType with optional tags.
func (r *GormProjectTypeRepository) Create(name string, tagIDs []uint) (*domain.ProjectType, error) {
	projectType := models.ProjectType{Name: name}

	if len(tagIDs) > 0 {
		var tags []models.Tag
		if err := r.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return nil, err
		}
		projectType.Tags = tags
	}

	if err := r.db.Create(&projectType).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&projectType), nil
}

// FindByID finds a ProjectType by its ID with tags preloaded.
func (r *GormProjectTypeRepository) FindByID(id uint) (*domain.ProjectType, error) {
	var projectType models.ProjectType
	if err := r.db.Preload("Tags").First(&projectType, id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&projectType), nil
}

// FindByName finds a ProjectType by its name with tags preloaded.
func (r *GormProjectTypeRepository) FindByName(name string) (*domain.ProjectType, error) {
	var projectType models.ProjectType
	if err := r.db.Preload("Tags").Where("name = ?", name).Find(&projectType).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&projectType), nil
}

// FindAllWithTags finds all ProjectTypes with tags preloaded.
func (r *GormProjectTypeRepository) FindAllWithTags() ([]domain.ProjectType, error) {
	var projectTypes []models.ProjectType
	if err := r.db.Preload("Tags").Find(&projectTypes).Error; err != nil {
		return nil, err
	}
	return r.toDomainSlice(projectTypes), nil
}

// GetAllNames returns all ProjectType names.
func (r *GormProjectTypeRepository) GetAllNames() []string {
	var projectTypes []models.ProjectType
	r.db.Find(&projectTypes)
	var names []string
	for _, pt := range projectTypes {
		names = append(names, pt.Name)
	}
	return names
}

// Update updates a ProjectType's name and tags.
func (r *GormProjectTypeRepository) Update(id uint, name string, tagIDs []uint) error {
	var projectType models.ProjectType
	if err := r.db.First(&projectType, id).Error; err != nil {
		return err
	}

	// Clear existing tags
	if err := r.db.Model(&projectType).Association("Tags").Clear(); err != nil {
		return err
	}

	projectType.Name = name

	if len(tagIDs) > 0 {
		var tags []models.Tag
		if err := r.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return err
		}
		projectType.Tags = tags
	}

	return r.db.Save(&projectType).Error
}

// Delete permanently deletes a ProjectType.
func (r *GormProjectTypeRepository) Delete(id uint) error {
	var projectType models.ProjectType
	if err := r.db.First(&projectType, id).Error; err != nil {
		return err
	}
	if err := r.db.Model(&projectType).Association("Tags").Clear(); err != nil {
		return err
	}
	return r.db.Unscoped().Delete(&projectType).Error
}

// toDomain converts a GORM model to a domain entity.
func (r *GormProjectTypeRepository) toDomain(m *models.ProjectType) *domain.ProjectType {
	d := &domain.ProjectType{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	for _, tag := range m.Tags {
		d.Tags = append(d.Tags, domain.Tag{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}
	return d
}

// toDomainSlice converts a slice of GORM models to domain entities.
func (r *GormProjectTypeRepository) toDomainSlice(ms []models.ProjectType) []domain.ProjectType {
	ds := make([]domain.ProjectType, len(ms))
	for i, m := range ms {
		ds[i] = *r.toDomain(&m)
	}
	return ds
}
