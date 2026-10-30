package note

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	domain "notes-service/internal/domain/note"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateNoteUseCase_Success(t *testing.T) {
	repo := new(mockRepository)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*note.Note")).Return(nil)
	uc := NewCreateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "Title", "hola", nil)

	assert.NoError(t, err)
	assert.Equal(t, "Title", n.Title)
	assert.Equal(t, "hola", n.Content)
	repo.AssertExpectations(t)
}

func TestCreateNoteUseCase_AcceptsEmptyTitle(t *testing.T) {
	repo := new(mockRepository)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*note.Note")).Return(nil)
	uc := NewCreateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "", "hola", nil)

	assert.NoError(t, err)
	assert.Equal(t, "", n.Title)
}

func TestCreateNoteUseCase_RejectsEmptyContent(t *testing.T) {
	repo := new(mockRepository)
	uc := NewCreateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "Title", "", nil)

	assert.Nil(t, n)
	assert.ErrorIs(t, err, domain.ErrEmptyContent)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateNoteUseCase_RepositoryError(t *testing.T) {
	repo := new(mockRepository)
	repoErr := errors.New("boom")
	repo.On("Create", mock.Anything, mock.AnythingOfType("*note.Note")).Return(repoErr)
	uc := NewCreateNoteUseCase(repo)

	_, err := uc.Execute(context.Background(), "Title", "hola", nil)

	assert.ErrorIs(t, err, repoErr)
}

func TestCreateNoteUseCase_PersistsNormalizedTags(t *testing.T) {
	repo := new(mockRepository)
	want := []string{"go", "backend"}
	repo.On("Create", mock.Anything, mock.MatchedBy(func(n *domain.Note) bool {
		return slices.Equal(n.Tags, want)
	})).Return(nil)
	uc := NewCreateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "Title", "hola", []string{"Go", " go ", "backend"})

	assert.NoError(t, err)
	assert.Equal(t, want, n.Tags)
	repo.AssertExpectations(t)
}

func TestCreateNoteUseCase_RejectsTagTooLong(t *testing.T) {
	repo := new(mockRepository)
	uc := NewCreateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "Title", "hola", []string{strings.Repeat("a", 21)})

	assert.Nil(t, n)
	assert.ErrorIs(t, err, domain.ErrTagTooLong)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateNoteUseCase_RejectsTooManyTags(t *testing.T) {
	repo := new(mockRepository)
	uc := NewCreateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "Title", "hola", []string{"a", "b", "c", "d", "e", "f"})

	assert.Nil(t, n)
	assert.ErrorIs(t, err, domain.ErrTooManyTags)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}
