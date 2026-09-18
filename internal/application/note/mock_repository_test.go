package note

import (
	"context"

	domain "notes-service/internal/domain/note"

	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, n *domain.Note) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

func (m *mockRepository) List(ctx context.Context) ([]*domain.Note, error) {
	args := m.Called(ctx)
	notes, _ := args.Get(0).([]*domain.Note)
	return notes, args.Error(1)
}

func (m *mockRepository) Get(ctx context.Context, id string) (*domain.Note, error) {
	args := m.Called(ctx, id)
	n, _ := args.Get(0).(*domain.Note)
	return n, args.Error(1)
}

func (m *mockRepository) Update(ctx context.Context, n *domain.Note) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
