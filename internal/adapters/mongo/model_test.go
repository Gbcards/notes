package mongoadapter

import (
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
