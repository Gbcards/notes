package note

import "context"

type DeleteNoteUseCase struct {
	repo Repository
}

func NewDeleteNoteUseCase(repo Repository) *DeleteNoteUseCase {
	return &DeleteNoteUseCase{repo: repo}
}

func (uc *DeleteNoteUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
