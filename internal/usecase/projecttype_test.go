package usecase

import (
	"testing"

	"github.com/niiharamegumu/chronowork/internal/repository/mock"
)

func TestProjectTypeUseCase_Create(t *testing.T) {
	tagRepo := mock.NewTagRepository()
	repo := mock.NewProjectTypeRepository(tagRepo)
	uc := NewProjectTypeUseCase(repo)

	pt, err := uc.Create("Test Project", nil)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if pt.Name != "Test Project" {
		t.Errorf("expected name 'Test Project', got '%s'", pt.Name)
	}
	if pt.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestProjectTypeUseCase_CreateWithTags(t *testing.T) {
	tagRepo := mock.NewTagRepository()
	tagUC := NewTagUseCase(tagRepo)
	tag1, _ := tagUC.Create("Tag1")
	tag2, _ := tagUC.Create("Tag2")

	repo := mock.NewProjectTypeRepository(tagRepo)
	uc := NewProjectTypeUseCase(repo)

	pt, err := uc.Create("Project With Tags", []uint{tag1.ID, tag2.ID})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if len(pt.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(pt.Tags))
	}
}

func TestProjectTypeUseCase_FindByID(t *testing.T) {
	tagRepo := mock.NewTagRepository()
	repo := mock.NewProjectTypeRepository(tagRepo)
	uc := NewProjectTypeUseCase(repo)

	created, _ := uc.Create("Test Project", nil)

	found, err := uc.FindByID(created.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Name != "Test Project" {
		t.Errorf("expected name 'Test Project', got '%s'", found.Name)
	}
}

func TestProjectTypeUseCase_FindByName(t *testing.T) {
	tagRepo := mock.NewTagRepository()
	repo := mock.NewProjectTypeRepository(tagRepo)
	uc := NewProjectTypeUseCase(repo)

	uc.Create("Find Me", nil)

	found, err := uc.FindByName("Find Me")
	if err != nil {
		t.Fatalf("FindByName failed: %v", err)
	}
	if found.Name != "Find Me" {
		t.Errorf("expected name 'Find Me', got '%s'", found.Name)
	}
}

func TestProjectTypeUseCase_GetAllNames(t *testing.T) {
	tagRepo := mock.NewTagRepository()
	repo := mock.NewProjectTypeRepository(tagRepo)
	uc := NewProjectTypeUseCase(repo)

	uc.Create("Project A", nil)
	uc.Create("Project B", nil)

	names := uc.GetAllNames()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestProjectTypeUseCase_Update(t *testing.T) {
	tagRepo := mock.NewTagRepository()
	repo := mock.NewProjectTypeRepository(tagRepo)
	uc := NewProjectTypeUseCase(repo)

	created, _ := uc.Create("Original", nil)

	err := uc.Update(created.ID, "Updated", nil)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := uc.FindByID(created.ID)
	if updated.Name != "Updated" {
		t.Errorf("expected name 'Updated', got '%s'", updated.Name)
	}
}

func TestProjectTypeUseCase_Delete(t *testing.T) {
	tagRepo := mock.NewTagRepository()
	repo := mock.NewProjectTypeRepository(tagRepo)
	uc := NewProjectTypeUseCase(repo)

	created, _ := uc.Create("To Delete", nil)

	err := uc.Delete(created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = uc.FindByID(created.ID)
	if err == nil {
		t.Error("expected error finding deleted project type")
	}
}
