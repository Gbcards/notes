package httpadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	domain "notes-service/internal/domain/note"

	"github.com/gin-gonic/gin"
)

type stubCreator struct {
	note *domain.Note
	err  error

	called  bool
	gotTags []string
}

func (s *stubCreator) Execute(ctx context.Context, title, content string, tags []string) (*domain.Note, error) {
	s.called = true
	s.gotTags = tags
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

	called  bool
	gotTags []string
}

func (s *stubUpdater) Execute(ctx context.Context, id, title, content string, tags []string) (*domain.Note, error) {
	s.called = true
	s.gotTags = tags
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

func doJSON(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// decodeObject decodes the body into a raw map so a JSON null is
// distinguishable from an empty array.
func decodeObject(t *testing.T, w *httptest.ResponseRecorder) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("decode response: %v: %s", err, w.Body.String())
	}
	return m
}

func assertRawTags(t *testing.T, raw json.RawMessage, want string) {
	t.Helper()
	if raw == nil {
		t.Fatalf("response has no tags key")
	}
	if got := string(raw); got != want {
		t.Fatalf("expected tags %s, got %s", want, got)
	}
}

func TestCreate_Success(t *testing.T) {
	n := &domain.Note{ID: "1", Title: "Title", Content: "hola", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	router := newTestRouter(&stubCreator{note: n}, nil, nil, nil)

	body, _ := json.Marshal(createNoteRequest{Title: "Title", Content: "hola"})
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreate_EmptyContentReturns400(t *testing.T) {
	router := newTestRouter(&stubCreator{err: domain.ErrEmptyContent}, nil, nil, nil)

	body, _ := json.Marshal(createNoteRequest{Title: "Title", Content: ""})
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreate_TitleTooLongReturns400(t *testing.T) {
	router := newTestRouter(&stubCreator{err: domain.ErrTitleTooLong}, nil, nil, nil)

	body, _ := json.Marshal(createNoteRequest{Title: "Title", Content: "hola"})
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreate_ForwardsTagsAndEchoesThem(t *testing.T) {
	n := &domain.Note{ID: "1", Content: "hola", Tags: []string{"go", "work"}, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	stub := &stubCreator{note: n}
	router := newTestRouter(stub, nil, nil, nil)

	w := doJSON(router, http.MethodPost, "/notes", `{"content":"hola","tags":[" Go ","work","go"]}`)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	if want := []string{" Go ", "work", "go"}; !reflect.DeepEqual(stub.gotTags, want) {
		t.Fatalf("expected stub to receive tags %q, got %q", want, stub.gotTags)
	}
	assertRawTags(t, decodeObject(t, w)["tags"], `["go","work"]`)
}

func TestCreate_NilTagsRenderAsEmptyArray(t *testing.T) {
	n := &domain.Note{ID: "1", Content: "hola", Tags: nil, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	router := newTestRouter(&stubCreator{note: n}, nil, nil, nil)

	w := doJSON(router, http.MethodPost, "/notes", `{"content":"hola"}`)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	assertRawTags(t, decodeObject(t, w)["tags"], `[]`)
}

func TestCreate_InvalidTagsTypeReturns400WithoutCallingUseCase(t *testing.T) {
	cases := map[string]string{
		"string":        `{"content":"hola","tags":"go"}`,
		"numeric array": `{"content":"hola","tags":[1, 2]}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			stub := &stubCreator{}
			router := newTestRouter(stub, nil, nil, nil)

			w := doJSON(router, http.MethodPost, "/notes", body)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
			}
			if stub.called {
				t.Fatalf("use case must not be called on malformed tags")
			}
		})
	}
}

func TestCreate_TagErrorsReturn400(t *testing.T) {
	for _, err := range []error{domain.ErrTagTooLong, domain.ErrTooManyTags} {
		t.Run(err.Error(), func(t *testing.T) {
			router := newTestRouter(&stubCreator{err: err}, nil, nil, nil)

			w := doJSON(router, http.MethodPost, "/notes", `{"content":"hola","tags":["go"]}`)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestList_Success(t *testing.T) {
	notes := []*domain.Note{{ID: "1", Title: "Title", Content: "a", CreatedAt: time.Now(), UpdatedAt: time.Now()}}
	router := newTestRouter(nil, stubLister{notes: notes}, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestList_NilTagsRenderAsEmptyArray(t *testing.T) {
	notes := []*domain.Note{
		{ID: "1", Content: "a", Tags: nil, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", Content: "b", Tags: []string{"go"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	router := newTestRouter(nil, stubLister{notes: notes}, nil, nil)

	w := doJSON(router, http.MethodGet, "/notes", "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var items []map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &items); err != nil {
		t.Fatalf("decode response: %v: %s", err, w.Body.String())
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(items))
	}
	assertRawTags(t, items[0]["tags"], `[]`)
	assertRawTags(t, items[1]["tags"], `["go"]`)
}

func TestUpdate_Success(t *testing.T) {
	n := &domain.Note{ID: "1", Title: "updated title", Content: "updated", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	router := newTestRouter(nil, nil, &stubUpdater{note: n}, nil)

	body, _ := json.Marshal(updateNoteRequest{Title: "updated title", Content: "updated"})
	req := httptest.NewRequest(http.MethodPut, "/notes/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_NotFoundReturns404(t *testing.T) {
	router := newTestRouter(nil, nil, &stubUpdater{err: domain.ErrNotFound}, nil)

	body, _ := json.Marshal(updateNoteRequest{Title: "updated title", Content: "updated"})
	req := httptest.NewRequest(http.MethodPut, "/notes/missing", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_TitleTooLongReturns400(t *testing.T) {
	router := newTestRouter(nil, nil, &stubUpdater{err: domain.ErrTitleTooLong}, nil)

	body, _ := json.Marshal(updateNoteRequest{Title: "updated title", Content: "updated"})
	req := httptest.NewRequest(http.MethodPut, "/notes/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_ForwardsTagsAndEchoesThem(t *testing.T) {
	n := &domain.Note{ID: "1", Content: "updated", Tags: []string{"go", "work"}, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	stub := &stubUpdater{note: n}
	router := newTestRouter(nil, nil, stub, nil)

	w := doJSON(router, http.MethodPut, "/notes/1", `{"content":"updated","tags":[" Go ","work","go"]}`)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if want := []string{" Go ", "work", "go"}; !reflect.DeepEqual(stub.gotTags, want) {
		t.Fatalf("expected stub to receive tags %q, got %q", want, stub.gotTags)
	}
	assertRawTags(t, decodeObject(t, w)["tags"], `["go","work"]`)
}

func TestUpdate_NilTagsRenderAsEmptyArray(t *testing.T) {
	n := &domain.Note{ID: "1", Content: "updated", Tags: nil, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	router := newTestRouter(nil, nil, &stubUpdater{note: n}, nil)

	w := doJSON(router, http.MethodPut, "/notes/1", `{"content":"updated"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	assertRawTags(t, decodeObject(t, w)["tags"], `[]`)
}

func TestUpdate_InvalidTagsTypeReturns400WithoutCallingUseCase(t *testing.T) {
	cases := map[string]string{
		"string":        `{"content":"updated","tags":"go"}`,
		"numeric array": `{"content":"updated","tags":[1, 2]}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			stub := &stubUpdater{}
			router := newTestRouter(nil, nil, stub, nil)

			w := doJSON(router, http.MethodPut, "/notes/1", body)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
			}
			if stub.called {
				t.Fatalf("use case must not be called on malformed tags")
			}
		})
	}
}

func TestUpdate_TagErrorsReturn400(t *testing.T) {
	for _, err := range []error{domain.ErrTagTooLong, domain.ErrTooManyTags} {
		t.Run(err.Error(), func(t *testing.T) {
			router := newTestRouter(nil, nil, &stubUpdater{err: err}, nil)

			w := doJSON(router, http.MethodPut, "/notes/1", `{"content":"updated","tags":["go"]}`)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
			}
		})
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
