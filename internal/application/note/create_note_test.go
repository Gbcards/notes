package note

import (
	"context"
	"errors"
	"testing"

	domain "notes-service/internal/domain/note"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateNoteUseCase_Success(t *testing.T) {
	repo := new(mockRepository)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*note.Note")).Return(nil)
	uc := NewCreateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "hola")

	assert.NoError(t, err)
	assert.Equal(t, "hola", n.Content)
	repo.AssertExpectations(t)
}

func TestCreateNoteUseCase_RejectsEmptyContent(t *testing.T) {
	repo := new(mockRepository)
	uc := NewCreateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "")

	assert.Nil(t, n)
	assert.ErrorIs(t, err, domain.ErrEmptyContent)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateNoteUseCase_RepositoryError(t *testing.T) {
	repo := new(mockRepository)
	repoErr := errors.New("boom")
	repo.On("Create", mock.Anything, mock.AnythingOfType("*note.Note")).Return(repoErr)
	uc := NewCreateNoteUseCase(repo)

	_, err := uc.Execute(context.Background(), "hola")

	assert.ErrorIs(t, err, repoErr)
}
