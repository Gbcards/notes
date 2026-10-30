package mongoadapter

import (
	"time"

	domain "notes-service/internal/domain/note"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type noteDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	Tags      []string           `bson:"tags"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func newDocument(n *domain.Note) noteDocument {
	return noteDocument{
		Title:     n.Title,
		Content:   n.Content,
		Tags:      n.Tags,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

func (d noteDocument) toDomain() *domain.Note {
	return &domain.Note{
		ID:        d.ID.Hex(),
		Title:     d.Title,
		Content:   d.Content,
		Tags:      d.Tags,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

// updateDocument builds the Mongo update for an existing note: a "$set" of the
// mutable fields. Tags are copied as-is; notes reaching here went through
// domain.New/Update, so they are never nil.
func updateDocument(n *domain.Note) bson.M {
	return bson.M{"$set": bson.M{
		"title":      n.Title,
		"content":    n.Content,
		"tags":       n.Tags,
		"updated_at": n.UpdatedAt,
	}}
}

// parseID converts an opaque note identifier into a Mongo ObjectID.
func parseID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
