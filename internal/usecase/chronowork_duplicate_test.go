package usecase

import (
	"testing"

	"github.com/niiharamegumu/chronowork/internal/repository/mock"
)

func TestChronoWorkUseCase_Create_DuplicateTitle(t *testing.T) {
	repo := mock.NewChronoWorkRepository()
	uc := NewChronoWorkUseCase(repo)

	// 1回目の作成（成功）
	_, err := uc.Create("Test Work", 1, 2)
	if err != nil {
		t.Fatalf("First create should succeed: %v", err)
	}

	// 2回目の作成（同じタイトル、今日作成なので失敗）
	_, err = uc.Create("Test Work", 1, 2)
	if err == nil {
		t.Error("Second create with same title should fail")
	}
	if err.Error() != "work with this title already exists today" {
		t.Errorf("Expected duplicate error, got: %v", err)
	}
}

func TestChronoWorkUseCase_Create_DifferentTitle(t *testing.T) {
	repo := mock.NewChronoWorkRepository()
	uc := NewChronoWorkUseCase(repo)

	// 異なるタイトルなら複数作成可能
	_, err := uc.Create("Work 1", 1, 2)
	if err != nil {
		t.Fatalf("Create Work 1 failed: %v", err)
	}

	_, err = uc.Create("Work 2", 1, 2)
	if err != nil {
		t.Fatalf("Create Work 2 failed: %v", err)
	}

	// 両方とも作成されていることを確認
	works, _ := uc.GetAll("", 0)
	if len(works) != 2 {
		t.Errorf("Expected 2 works, got %d", len(works))
	}
}
