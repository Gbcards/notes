package note

import (
	"context"
	"slices"
	"strings"
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

	n, err := uc.Execute(context.Background(), "1", "new title", "new", nil)

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

	n, err := uc.Execute(context.Background(), "1", "", "new", nil)

	assert.NoError(t, err)
	assert.Equal(t, "", n.Title)
}

func TestUpdateNoteUseCase_NotFound(t *testing.T) {
	repo := new(mockRepository)
	repo.On("Get", mock.Anything, "missing").Return(nil, domain.ErrNotFound)
	uc := NewUpdateNoteUseCase(repo)

	_, err := uc.Execute(context.Background(), "missing", "title", "new", nil)

	assert.ErrorIs(t, err, domain.ErrNotFound)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestUpdateNoteUseCase_RejectsEmptyContent(t *testing.T) {
	repo := new(mockRepository)
	created := time.Now()
	existing := &domain.Note{ID: "1", Title: "old title", Content: "old", CreatedAt: created, UpdatedAt: created}
	repo.On("Get", mock.Anything, "1").Return(existing, nil)
	uc := NewUpdateNoteUseCase(repo)

	_, err := uc.Execute(context.Background(), "1", "new title", "", nil)

	assert.ErrorIs(t, err, domain.ErrEmptyContent)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestUpdateNoteUseCase_PersistsNormalizedTags(t *testing.T) {
	repo := new(mockRepository)
	created := time.Now().Add(-time.Hour)
	existing := &domain.Note{ID: "1", Content: "old", Tags: []string{"old"}, CreatedAt: created, UpdatedAt: created}
	want := []string{"go", "backend"}
	repo.On("Get", mock.Anything, "1").Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(n *domain.Note) bool {
		return slices.Equal(n.Tags, want)
	})).Return(nil)
	uc := NewUpdateNoteUseCase(repo)

	n, err := uc.Execute(context.Background(), "1", "title", "new", []string{"Go", " go ", "backend"})

	assert.NoError(t, err)
	assert.Equal(t, want, n.Tags)
	repo.AssertExpectations(t)
}

func TestUpdateNoteUseCase_RejectsTagTooLong(t *testing.T) {
	repo := new(mockRepository)
	created := time.Now()
	existing := &domain.Note{ID: "1", Content: "old", CreatedAt: created, UpdatedAt: created}
	repo.On("Get", mock.Anything, "1").Return(existing, nil)
	uc := NewUpdateNoteUseCase(repo)

	_, err := uc.Execute(context.Background(), "1", "title", "new", []string{strings.Repeat("a", 21)})

	assert.ErrorIs(t, err, domain.ErrTagTooLong)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestUpdateNoteUseCase_RejectsTooManyTags(t *testing.T) {
	repo := new(mockRepository)
	created := time.Now()
	existing := &domain.Note{ID: "1", Content: "old", CreatedAt: created, UpdatedAt: created}
	repo.On("Get", mock.Anything, "1").Return(existing, nil)
	uc := NewUpdateNoteUseCase(repo)

	_, err := uc.Execute(context.Background(), "1", "title", "new", []string{"a", "b", "c", "d", "e", "f"})

	assert.ErrorIs(t, err, domain.ErrTooManyTags)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}
