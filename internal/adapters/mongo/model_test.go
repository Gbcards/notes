package mongoadapter

import (
	"testing"

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
