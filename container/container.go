// Package container provides dependency injection container for the application.
package container

import (
	"github.com/niiharamegumu/chronowork/internal/repository"
	"github.com/niiharamegumu/chronowork/internal/usecase"
	"gorm.io/gorm"
)

// Container holds all application dependencies.
type Container struct {
	DB *gorm.DB

	// Repositories
	ChronoWorkRepo  repository.ChronoWorkRepository
	TagRepo         repository.TagRepository
	ProjectTypeRepo repository.ProjectTypeRepository
	SettingRepo     repository.SettingRepository

	// Use Cases
	ChronoWorkUC  *usecase.ChronoWorkUseCase
	TagUC         *usecase.TagUseCase
	ProjectTypeUC *usecase.ProjectTypeUseCase
	SettingUC     *usecase.SettingUseCase
}

// New creates a new Container with all dependencies initialized.
func New(db *gorm.DB) *Container {
	// Initialize repositories
	chronoWorkRepo := repository.NewGormChronoWorkRepository(db)
	tagRepo := repository.NewGormTagRepository(db)
	projectTypeRepo := repository.NewGormProjectTypeRepository(db)
	settingRepo := repository.NewGormSettingRepository(db)

	// Initialize use cases
	chronoWorkUC := usecase.NewChronoWorkUseCase(chronoWorkRepo)
	tagUC := usecase.NewTagUseCase(tagRepo)
	projectTypeUC := usecase.NewProjectTypeUseCase(projectTypeRepo)
	settingUC := usecase.NewSettingUseCase(settingRepo)

	return &Container{
		DB: db,

		ChronoWorkRepo:  chronoWorkRepo,
		TagRepo:         tagRepo,
		ProjectTypeRepo: projectTypeRepo,
		SettingRepo:     settingRepo,

		ChronoWorkUC:  chronoWorkUC,
		TagUC:         tagUC,
		ProjectTypeUC: projectTypeUC,
		SettingUC:     settingUC,
	}
}
