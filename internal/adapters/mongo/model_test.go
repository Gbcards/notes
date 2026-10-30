package mongoadapter

import (
	"reflect"
	"testing"
	"time"

	domain "notes-service/internal/domain/note"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestParseID_RoundTripsValidHex(t *testing.T) {
	original := primitive.NewObjectID()

	parsed, err := parseID(original.Hex())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if parsed != original {
		t.Errorf("expected %v, got %v", original, parsed)
	}
}

func TestParseID_RejectsMalformedHex(t *testing.T) {
	_, err := parseID("not-a-valid-object-id")
	if err == nil {
		t.Fatal("expected an error for malformed hex, got nil")
	}
}

func TestNewDocument_IncludesTitle(t *testing.T) {
	now := time.Now()
	n := &domain.Note{Title: "Title", Content: "hola", CreatedAt: now, UpdatedAt: now}

	doc := newDocument(n)

	if doc.Title != "Title" {
		t.Errorf("expected title %q, got %q", "Title", doc.Title)
	}
}

func TestToDomain_IncludesTitle(t *testing.T) {
	doc := noteDocument{ID: primitive.NewObjectID(), Title: "Title", Content: "hola"}

	n := doc.toDomain()

	if n.Title != "Title" {
		t.Errorf("expected title %q, got %q", "Title", n.Title)
	}
}

func TestNoteDocument_MissingTitleKeyDecodesToEmptyTitle(t *testing.T) {
	raw, err := bson.Marshal(bson.M{
		"_id":        primitive.NewObjectID(),
		"content":    "hola",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}

	var doc noteDocument
	if err := bson.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}

	if doc.Title != "" {
		t.Errorf("expected empty title for document without a stored title key, got %q", doc.Title)
	}
}

func TestNewDocument_IncludesTags(t *testing.T) {
	now := time.Now()
	n := &domain.Note{Content: "hola", Tags: []string{"go", "mongo"}, CreatedAt: now, UpdatedAt: now}

	doc := newDocument(n)

	if !reflect.DeepEqual(doc.Tags, []string{"go", "mongo"}) {
		t.Errorf("expected tags %v, got %v", []string{"go", "mongo"}, doc.Tags)
	}
}

func TestToDomain_IncludesTags(t *testing.T) {
	doc := noteDocument{ID: primitive.NewObjectID(), Content: "hola", Tags: []string{"go", "mongo"}}

	n := doc.toDomain()

	if !reflect.DeepEqual(n.Tags, []string{"go", "mongo"}) {
		t.Errorf("expected tags %v, got %v", []string{"go", "mongo"}, n.Tags)
	}
}

func TestNoteDocument_MissingTagsKeyDecodesToNilTags(t *testing.T) {
	raw, err := bson.Marshal(bson.M{
		"_id":        primitive.NewObjectID(),
		"title":      "Title",
		"content":    "hola",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}

	var doc noteDocument
	if err := bson.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}

	n := doc.toDomain()
	if n.Tags != nil {
		t.Errorf("expected nil tags for document without a stored tags key, got %#v", n.Tags)
	}
}

func TestUpdateDocument_SetsAllMutableFields(t *testing.T) {
	updatedAt := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
	n := &domain.Note{
		ID:        primitive.NewObjectID().Hex(),
		Title:     "Title",
		Content:   "hola",
		Tags:      []string{"go", "mongo"},
		CreatedAt: updatedAt.Add(-time.Hour),
		UpdatedAt: updatedAt,
	}

	update := updateDocument(n)

	if len(update) != 1 {
		t.Fatalf("expected only a $set operator, got %v", update)
	}
	set, ok := update["$set"].(bson.M)
	if !ok {
		t.Fatalf("expected $set to be a bson.M, got %T", update["$set"])
	}

	want := bson.M{
		"title":      "Title",
		"content":    "hola",
		"tags":       []string{"go", "mongo"},
		"updated_at": updatedAt,
	}
	if !reflect.DeepEqual(set, want) {
		t.Errorf("expected $set %v, got %v", want, set)
	}
}
