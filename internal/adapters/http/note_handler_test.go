package httpadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "notes-service/internal/domain/note"

	"github.com/gin-gonic/gin"
)

type stubCreator struct {
	note *domain.Note
	err  error
}

func (s stubCreator) Execute(ctx context.Context, content string) (*domain.Note, error) {
	return s.note, s.err
}

type stubLister struct {
	notes []*domain.Note
	err   error
}

func (s stubLister) Execute(ctx context.Context) ([]*domain.Note, error) {
	return s.notes, s.err
}

type stubUpdater struct {
	note *domain.Note
	err  error
}

func (s stubUpdater) Execute(ctx context.Context, id, content string) (*domain.Note, error) {
	return s.note, s.err
}

type stubDeleter struct {
	err error
}

func (s stubDeleter) Execute(ctx context.Context, id string) error {
	return s.err
}

func newTestRouter(create NoteCreator, list NoteLister, update NoteUpdater, del NoteDeleter) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter(NewNoteHandler(create, list, update, del))
}

func TestCreate_Success(t *testing.T) {
	n := &domain.Note{ID: "1", Content: "hola", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	router := newTestRouter(stubCreator{note: n}, nil, nil, nil)

	body, _ := json.Marshal(createNoteRequest{Content: "hola"})
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreate_EmptyContentReturns400(t *testing.T) {
	router := newTestRouter(stubCreator{err: domain.ErrEmptyContent}, nil, nil, nil)

	body, _ := json.Marshal(createNoteRequest{Content: ""})
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestList_Success(t *testing.T) {
	notes := []*domain.Note{{ID: "1", Content: "a", CreatedAt: time.Now(), UpdatedAt: time.Now()}}
	router := newTestRouter(nil, stubLister{notes: notes}, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_Success(t *testing.T) {
	n := &domain.Note{ID: "1", Content: "updated", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	router := newTestRouter(nil, nil, stubUpdater{note: n}, nil)

	body, _ := json.Marshal(updateNoteRequest{Content: "updated"})
	req := httptest.NewRequest(http.MethodPut, "/notes/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_NotFoundReturns404(t *testing.T) {
	router := newTestRouter(nil, nil, stubUpdater{err: domain.ErrNotFound}, nil)

	body, _ := json.Marshal(updateNoteRequest{Content: "updated"})
	req := httptest.NewRequest(http.MethodPut, "/notes/missing", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDelete_Success(t *testing.T) {
	router := newTestRouter(nil, nil, nil, stubDeleter{})

	req := httptest.NewRequest(http.MethodDelete, "/notes/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDelete_NotFoundReturns404(t *testing.T) {
	router := newTestRouter(nil, nil, nil, stubDeleter{err: domain.ErrNotFound})

	req := httptest.NewRequest(http.MethodDelete, "/notes/missing", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}
