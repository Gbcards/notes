package note

import (
	"context"
	"time"

	domain "notes-service/internal/domain/note"
)

type UpdateNoteUseCase struct {
	repo Repository
}

func NewUpdateNoteUseCase(repo Repository) *UpdateNoteUseCase {
	return &UpdateNoteUseCase{repo: repo}
}

func (uc *UpdateNoteUseCase) Execute(ctx context.Context, id, title, content string) (*domain.Note, error) {
	n, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := n.Update(title, content, time.Now()); err != nil {
		return nil, err
	}
	if err := uc.repo.Update(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}
