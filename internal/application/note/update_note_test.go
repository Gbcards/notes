package note

import (
	"context"
	"testing"
	"time"

	domain "notes-service/internal/domain/note"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateNoteUseCase_Success(t *testing.T) {
	repo := new(mockRepository)
	created := time.Now().Add(-time.Hour)
	existing := &domain.Note{ID: "1", Title: "old title", Content: "old", CreatedAt: created, UpdatedAt: created}
	repo.On("Get", mock.Anything, "1").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*note.Note")).Return(nil)
	uc := NewUpdateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "1", "new title", "new")

	assert.NoError(t, err)
	assert.Equal(t, "new title", n.Title)
	assert.Equal(t, "new", n.Content)
	assert.True(t, n.UpdatedAt.After(created))
	assert.True(t, n.CreatedAt.Equal(created))
}

func TestUpdateNoteUseCase_AcceptsEmptyTitle(t *testing.T) {
	repo := new(mockRepository)
	created := time.Now().Add(-time.Hour)
	existing := &domain.Note{ID: "1", Title: "old title", Content: "old", CreatedAt: created, UpdatedAt: created}
	repo.On("Get", mock.Anything, "1").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*note.Note")).Return(nil)
	uc := NewUpdateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "1", "", "new")

	assert.NoError(t, err)
	assert.Equal(t, "", n.Title)
}

func TestUpdateNoteUseCase_NotFound(t *testing.T) {
	repo := new(mockRepository)
	repo.On("Get", mock.Anything, "missing").Return(nil, domain.ErrNotFound)
	uc := NewUpdateNoteUseCase(repo)

	_, err := uc.Execute(context.Background(), "missing", "title", "new")

	assert.ErrorIs(t, err, domain.ErrNotFound)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestUpdateNoteUseCase_RejectsEmptyContent(t *testing.T) {
	repo := new(mockRepository)
	created := time.Now()
	existing := &domain.Note{ID: "1", Title: "old title", Content: "old", CreatedAt: created, UpdatedAt: created}
	repo.On("Get", mock.Anything, "1").Return(existing, nil)
	uc := NewUpdateNoteUseCase(repo)

	_, err := uc.Execute(context.Background(), "1", "new title", "")

	assert.ErrorIs(t, err, domain.ErrEmptyContent)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}
