package note

import (
	"context"
	"time"

	domain "notes-service/internal/domain/note"
)

type CreateNoteUseCase struct {
	repo Repository
}

func NewCreateNoteUseCase(repo Repository) *CreateNoteUseCase {
	return &CreateNoteUseCase{repo: repo}
}

func (uc *CreateNoteUseCase) Execute(ctx context.Context, content string) (*domain.Note, error) {
	n, err := domain.New(content, time.Now())
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}
