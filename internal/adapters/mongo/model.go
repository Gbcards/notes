package mongoadapter

import (
	"time"

	domain "notes-service/internal/domain/note"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type noteDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func newDocument(n *domain.Note) noteDocument {
	return noteDocument{
		Title:     n.Title,
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

func (d noteDocument) toDomain() *domain.Note {
	return &domain.Note{
		ID:        d.ID.Hex(),
		Title:     d.Title,
		Content:   d.Content,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

// parseID converts an opaque note identifier into a Mongo ObjectID.
func parseID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
