package mongoadapter

import (
	"context"
	"errors"

	application "notes-service/internal/application/note"
	domain "notes-service/internal/domain/note"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NoteRepository struct {
	collection *mongo.Collection
}

func NewNoteRepository(collection *mongo.Collection) *NoteRepository {
	return &NoteRepository{collection: collection}
}

var _ application.Repository = (*NoteRepository)(nil)

func (r *NoteRepository) Create(ctx context.Context, n *domain.Note) error {
	doc := newDocument(n)
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	n.ID = result.InsertedID.(interface{ Hex() string }).Hex()
	return nil
}

func (r *NoteRepository) List(ctx context.Context) ([]*domain.Note, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	notes := []*domain.Note{}
	for cursor.Next(ctx) {
		var doc noteDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		notes = append(notes, doc.toDomain())
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *NoteRepository) Get(ctx context.Context, id string) (*domain.Note, error) {
	objectID, err := parseID(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	var doc noteDocument
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return doc.toDomain(), nil
}

func (r *NoteRepository) Update(ctx context.Context, n *domain.Note) error {
	objectID, err := parseID(n.ID)
	if err != nil {
		return domain.ErrNotFound
	}

	update := bson.M{"$set": bson.M{
		"content":    n.Content,
		"updated_at": n.UpdatedAt,
	}}
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *NoteRepository) Delete(ctx context.Context, id string) error {
	objectID, err := parseID(id)
	if err != nil {
		return domain.ErrNotFound
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}
