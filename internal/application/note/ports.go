package note

import (
	"context"

	domain "notes-service/internal/domain/note"
)

// Repository is the outbound port implemented by persistence adapters.
type Repository interface {
	// Create persists n and sets n.ID to the assigned identifier.
	Create(ctx context.Context, n *domain.Note) error
	List(ctx context.Context) ([]*domain.Note, error)
	Get(ctx context.Context, id string) (*domain.Note, error)
	Update(ctx context.Context, n *domain.Note) error
	Delete(ctx context.Context, id string) error
}
