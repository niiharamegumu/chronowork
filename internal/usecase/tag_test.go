package usecase

import (
	"testing"

	"github.com/niiharamegumu/chronowork/internal/repository/mock"
)

func TestTagUseCase_Create(t *testing.T) {
	repo := mock.NewTagRepository()
	uc := NewTagUseCase(repo)

	tag, err := uc.Create("Test Tag")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if tag.Name != "Test Tag" {
		t.Errorf("expected name 'Test Tag', got '%s'", tag.Name)
	}
	if tag.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestTagUseCase_FindByID(t *testing.T) {
	repo := mock.NewTagRepository()
	uc := NewTagUseCase(repo)

	created, _ := uc.Create("Test Tag")

	found, err := uc.FindByID(created.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Name != "Test Tag" {
		t.Errorf("expected name 'Test Tag', got '%s'", found.Name)
	}
}

func TestTagUseCase_FindAll(t *testing.T) {
	repo := mock.NewTagRepository()
	uc := NewTagUseCase(repo)

	uc.Create("Tag 1")
	uc.Create("Tag 2")

	tags, err := uc.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
}

func TestTagUseCase_GetAllNames(t *testing.T) {
	repo := mock.NewTagRepository()
	uc := NewTagUseCase(repo)

	uc.Create("Alpha")
	uc.Create("Beta")

	names := uc.GetAllNames()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestTagUseCase_Update(t *testing.T) {
	repo := mock.NewTagRepository()
	uc := NewTagUseCase(repo)

	created, _ := uc.Create("Original")

	err := uc.Update(created.ID, "Updated")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := uc.FindByID(created.ID)
	if updated.Name != "Updated" {
		t.Errorf("expected name 'Updated', got '%s'", updated.Name)
	}
}

func TestTagUseCase_Delete(t *testing.T) {
	repo := mock.NewTagRepository()
	uc := NewTagUseCase(repo)

	created, _ := uc.Create("To Delete")

	err := uc.Delete(created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = uc.FindByID(created.ID)
	if err == nil {
		t.Error("expected error finding deleted tag")
	}
}
