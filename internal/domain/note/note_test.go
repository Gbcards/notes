package note

import (
	"errors"
	"testing"
	"time"
)

func TestNew_AcceptsNonEmptyContent(t *testing.T) {
	now := time.Now()

	n, err := New("Title", "hola", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Title != "Title" {
		t.Errorf("expected title %q, got %q", "Title", n.Title)
	}
	if n.Content != "hola" {
		t.Errorf("expected content %q, got %q", "hola", n.Content)
	}
	if !n.CreatedAt.Equal(now) || !n.UpdatedAt.Equal(now) {
		t.Errorf("expected CreatedAt and UpdatedAt to equal %v, got %v and %v", now, n.CreatedAt, n.UpdatedAt)
	}
}

func TestNew_AcceptsEmptyTitle(t *testing.T) {
	n, err := New("", "hola", time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Title != "" {
		t.Errorf("expected empty title, got %q", n.Title)
	}
}

func TestNew_RejectsEmptyContent(t *testing.T) {
	_, err := New("Title", "", time.Now())
	if !errors.Is(err, ErrEmptyContent) {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
}

func TestUpdate_Success(t *testing.T) {
	created := time.Now().Add(-time.Hour)
	n, err := New("original title", "original", created)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedAt := time.Now()
	if err := n.Update("updated title", "updated", updatedAt); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Title != "updated title" {
		t.Errorf("expected title %q, got %q", "updated title", n.Title)
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

func TestUpdate_RejectsEmptyContent(t *testing.T) {
	n, err := New("original title", "original", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := n.Update("new title", "", time.Now()); !errors.Is(err, ErrEmptyContent) {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
	if n.Content != "original" {
		t.Errorf("expected content to remain %q, got %q", "original", n.Content)
	}
	if n.Title != "original title" {
		t.Errorf("expected title to remain %q, got %q", "original title", n.Title)
	}
}
