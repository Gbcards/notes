package note

import (
	"errors"
	"testing"
	"time"
)

func TestNew_AcceptsNonEmptyContent(t *testing.T) {
	now := time.Now()

	n, err := New("hola", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Content != "hola" {
		t.Errorf("expected content %q, got %q", "hola", n.Content)
	}
	if !n.CreatedAt.Equal(now) || !n.UpdatedAt.Equal(now) {
		t.Errorf("expected CreatedAt and UpdatedAt to equal %v, got %v and %v", now, n.CreatedAt, n.UpdatedAt)
	}
}

func TestNew_RejectsEmptyContent(t *testing.T) {
	_, err := New("", time.Now())
	if !errors.Is(err, ErrEmptyContent) {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
}

func TestUpdateContent_Success(t *testing.T) {
	created := time.Now().Add(-time.Hour)
	n, err := New("original", created)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedAt := time.Now()
	if err := n.UpdateContent("updated", updatedAt); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Content != "updated" {
		t.Errorf("expected content %q, got %q", "updated", n.Content)
	}
	if !n.UpdatedAt.Equal(updatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", updatedAt, n.UpdatedAt)
	}
	if !n.CreatedAt.Equal(created) {
		t.Errorf("expected CreatedAt to remain %v, got %v", created, n.CreatedAt)
	}
}

func TestUpdateContent_RejectsEmptyContent(t *testing.T) {
	n, err := New("original", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := n.UpdateContent("", time.Now()); !errors.Is(err, ErrEmptyContent) {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
	if n.Content != "original" {
		t.Errorf("expected content to remain %q, got %q", "original", n.Content)
	}
}
