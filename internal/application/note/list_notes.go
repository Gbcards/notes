package note

import (
	"context"

	domain "notes-service/internal/domain/note"
)

type ListNotesUseCase struct {
	repo Repository
}

func NewListNotesUseCase(repo Repository) *ListNotesUseCase {
	return &ListNotesUseCase{repo: repo}
}

func (uc *ListNotesUseCase) Execute(ctx context.Context) ([]*domain.Note, error) {
	return uc.repo.List(ctx)
}
