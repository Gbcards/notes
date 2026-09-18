package note

import (
	"context"
	"testing"
	"time"

	domain "notes-service/internal/domain/note"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListNotesUseCase_ReturnsAllNotes(t *testing.T) {
	repo := new(mockRepository)
	expected := []*domain.Note{
		{ID: "1", Content: "a", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", Content: "b", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	repo.On("List", mock.Anything).Return(expected, nil)
	uc := NewListNotesUseCase(repo)

	notes, err := uc.Execute(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, notes)
}

func TestListNotesUseCase_ReturnsEmptyListWhenNoneExist(t *testing.T) {
	repo := new(mockRepository)
	repo.On("List", mock.Anything).Return([]*domain.Note{}, nil)
	uc := NewListNotesUseCase(repo)

	notes, err := uc.Execute(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, notes)
}
