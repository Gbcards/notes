package note

import (
	"context"
	"testing"

	domain "notes-service/internal/domain/note"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteNoteUseCase_Success(t *testing.T) {
	repo := new(mockRepository)
	repo.On("Delete", mock.Anything, "1").Return(nil)
	uc := NewDeleteNoteUseCase(repo)

	err := uc.Execute(context.Background(), "1")

	assert.NoError(t, err)
}

func TestDeleteNoteUseCase_NotFound(t *testing.T) {
	repo := new(mockRepository)
	repo.On("Delete", mock.Anything, "missing").Return(domain.ErrNotFound)
	uc := NewDeleteNoteUseCase(repo)

	err := uc.Execute(context.Background(), "missing")

	assert.ErrorIs(t, err, domain.ErrNotFound)
}
