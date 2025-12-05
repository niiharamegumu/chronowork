package usecase

import (
	"testing"
	"time"

	"github.com/niiharamegumu/chronowork/internal/repository/mock"
)

func TestChronoWorkUseCase_Create(t *testing.T) {
	repo := mock.NewChronoWorkRepository()
	uc := NewChronoWorkUseCase(repo)

	// Test creation
	cw, err := uc.Create("Test Work", 1, 2)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if cw.Title != "Test Work" {
		t.Errorf("expected title 'Test Work', got '%s'", cw.Title)
	}
	if cw.ProjectTypeID != 1 {
		t.Errorf("expected ProjectTypeID 1, got %d", cw.ProjectTypeID)
	}
	if cw.TagID != 2 {
		t.Errorf("expected TagID 2, got %d", cw.TagID)
	}
	if cw.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestChronoWorkUseCase_FindByID(t *testing.T) {
	repo := mock.NewChronoWorkRepository()
	uc := NewChronoWorkUseCase(repo)

	// Create a work entry
	created, _ := uc.Create("Test Work", 0, 0)

	// Find it
	found, err := uc.FindByID(created.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Title != "Test Work" {
		t.Errorf("expected title 'Test Work', got '%s'", found.Title)
	}

	// Find non-existent
	_, err = uc.FindByID(999)
	if err == nil {
		t.Error("expected error for non-existent ID")
	}
}

func TestChronoWorkUseCase_Update(t *testing.T) {
	repo := mock.NewChronoWorkRepository()
	uc := NewChronoWorkUseCase(repo)

	// Create
	created, _ := uc.Create("Original", 0, 0)

	// Update
	err := uc.Update(created.ID, "Updated", 1, 2)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify
	updated, _ := uc.FindByID(created.ID)
	if updated.Title != "Updated" {
		t.Errorf("expected title 'Updated', got '%s'", updated.Title)
	}
	if updated.ProjectTypeID != 1 {
		t.Errorf("expected ProjectTypeID 1, got %d", updated.ProjectTypeID)
	}
}

func TestChronoWorkUseCase_Delete(t *testing.T) {
	repo := mock.NewChronoWorkRepository()
	uc := NewChronoWorkUseCase(repo)

	// Create
	created, _ := uc.Create("To Delete", 0, 0)

	// Delete
	err := uc.Delete(created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	_, err = uc.FindByID(created.ID)
	if err == nil {
		t.Error("expected error finding deleted record")
	}
}

func TestChronoWorkUseCase_Tracking(t *testing.T) {
	repo := mock.NewChronoWorkRepository()
	uc := NewChronoWorkUseCase(repo)

	// Create
	created, _ := uc.Create("Track Test", 0, 0)

	// Initial state should not be tracking
	tracking, _ := uc.FindTracking()
	if len(tracking) != 0 {
		t.Error("expected no tracking entries initially")
	}

	// Start tracking
	err := uc.StartTracking(created.ID)
	if err != nil {
		t.Fatalf("StartTracking failed: %v", err)
	}

	// Verify tracking
	tracking, _ = uc.FindTracking()
	if len(tracking) != 1 {
		t.Errorf("expected 1 tracking entry, got %d", len(tracking))
	}

	// Verify StartTime is set
	found, _ := uc.FindByID(created.ID)
	if found.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}

	// Stop tracking
	time.Sleep(10 * time.Millisecond) // Small delay for time calculation
	err = uc.StopTracking(created.ID)
	if err != nil {
		t.Fatalf("StopTracking failed: %v", err)
	}

	// Verify stopped
	tracking, _ = uc.FindTracking()
	if len(tracking) != 0 {
		t.Errorf("expected 0 tracking entries after stop, got %d", len(tracking))
	}
}

func TestChronoWorkUseCase_UpdateConfirmed(t *testing.T) {
	repo := mock.NewChronoWorkRepository()
	uc := NewChronoWorkUseCase(repo)

	// Create
	created, _ := uc.Create("Confirm Test", 0, 0)

	// Initially not confirmed
	found, _ := uc.FindByID(created.ID)
	if found.Confirmed {
		t.Error("expected Confirmed to be false initially")
	}

	// Confirm
	err := uc.UpdateConfirmed(created.ID, true)
	if err != nil {
		t.Fatalf("UpdateConfirmed failed: %v", err)
	}

	// Verify
	found, _ = uc.FindByID(created.ID)
	if !found.Confirmed {
		t.Error("expected Confirmed to be true")
	}
}

func TestChronoWorkUseCase_UpdateTotalSeconds(t *testing.T) {
	repo := mock.NewChronoWorkRepository()
	uc := NewChronoWorkUseCase(repo)

	// Create
	created, _ := uc.Create("Timer Test", 0, 0)

	// Update total seconds
	err := uc.UpdateTotalSeconds(created.ID, 3600)
	if err != nil {
		t.Fatalf("UpdateTotalSeconds failed: %v", err)
	}

	// Verify
	found, _ := uc.FindByID(created.ID)
	if found.TotalSeconds != 3600 {
		t.Errorf("expected TotalSeconds 3600, got %d", found.TotalSeconds)
	}
}
