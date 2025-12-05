package usecase

import (
	"testing"

	"github.com/niiharamegumu/chronowork/internal/domain"
	"github.com/niiharamegumu/chronowork/internal/repository/mock"
)

func TestSettingUseCase_Get(t *testing.T) {
	repo := mock.NewSettingRepository()
	uc := NewSettingUseCase(repo)

	setting, err := uc.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Default values
	if setting.PersonDay != 8 {
		t.Errorf("expected default PersonDay 8, got %d", setting.PersonDay)
	}
	if setting.RelativeDate != 0 {
		t.Errorf("expected default RelativeDate 0, got %d", setting.RelativeDate)
	}
	if !setting.DisplayAsPersonDay {
		t.Error("expected default DisplayAsPersonDay true")
	}
}

func TestSettingUseCase_Update(t *testing.T) {
	repo := mock.NewSettingRepository()
	uc := NewSettingUseCase(repo)

	// Get default
	setting, _ := uc.Get()

	// Update
	updated := &domain.Setting{
		ID:                 setting.ID,
		RelativeDate:       7,
		PersonDay:          6,
		DisplayAsPersonDay: false,
		DownloadPath:       "/tmp/export",
	}
	err := uc.Update(updated)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify
	result, _ := uc.Get()
	if result.RelativeDate != 7 {
		t.Errorf("expected RelativeDate 7, got %d", result.RelativeDate)
	}
	if result.PersonDay != 6 {
		t.Errorf("expected PersonDay 6, got %d", result.PersonDay)
	}
	if result.DisplayAsPersonDay {
		t.Error("expected DisplayAsPersonDay false")
	}
	if result.DownloadPath != "/tmp/export" {
		t.Errorf("expected DownloadPath '/tmp/export', got '%s'", result.DownloadPath)
	}
}
