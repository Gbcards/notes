package note

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestNew_AcceptsNonEmptyContent(t *testing.T) {
	now := time.Now()

	n, err := New("Title", "hola", nil, now)
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
	n, err := New("", "hola", nil, time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Title != "" {
		t.Errorf("expected empty title, got %q", n.Title)
	}
}

func TestNew_RejectsEmptyContent(t *testing.T) {
	_, err := New("Title", "", nil, time.Now())
	if !errors.Is(err, ErrEmptyContent) {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
}

func TestNew_AcceptsTitleAtMaxLength(t *testing.T) {
	title := strings.Repeat("a", 200)

	n, err := New(title, "hola", nil, time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Title != title {
		t.Errorf("expected title to be stored unchanged, got %q", n.Title)
	}
}

func TestNew_RejectsTitleExceedingMaxLength(t *testing.T) {
	title := strings.Repeat("a", 201)

	_, err := New(title, "hola", nil, time.Now())
	if !errors.Is(err, ErrTitleTooLong) {
		t.Fatalf("expected ErrTitleTooLong, got %v", err)
	}
}

func TestNew_AcceptsMultiByteTitleAtMaxRuneLength(t *testing.T) {
	// "ñ" and "á" are 2 bytes each in UTF-8, so 200 of them is 400 bytes
	// but only 200 runes -- this proves rune-counting, not byte-length, is used.
	title := strings.Repeat("ñá", 100)

	n, err := New(title, "hola", nil, time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Title != title {
		t.Errorf("expected title to be stored unchanged, got %q", n.Title)
	}
}

func TestUpdate_Success(t *testing.T) {
	created := time.Now().Add(-time.Hour)
	n, err := New("original title", "original", nil, created)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedAt := time.Now()
	if err := n.Update("updated title", "updated", nil, updatedAt); err != nil {
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
	n, err := New("original title", "original", nil, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := n.Update("new title", "", nil, time.Now()); !errors.Is(err, ErrEmptyContent) {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
	if n.Content != "original" {
		t.Errorf("expected content to remain %q, got %q", "original", n.Content)
	}
	if n.Title != "original title" {
		t.Errorf("expected title to remain %q, got %q", "original title", n.Title)
	}
}

func TestUpdate_AcceptsTitleAtMaxLength(t *testing.T) {
	n, err := New("original title", "original", nil, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	title := strings.Repeat("a", 200)
	if err := n.Update(title, "updated", nil, time.Now()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Title != title {
		t.Errorf("expected title to be stored unchanged, got %q", n.Title)
	}
}

func TestUpdate_RejectsTitleExceedingMaxLength(t *testing.T) {
	n, err := New("original title", "original", nil, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	title := strings.Repeat("a", 201)
	if err := n.Update(title, "updated", nil, time.Now()); !errors.Is(err, ErrTitleTooLong) {
		t.Fatalf("expected ErrTitleTooLong, got %v", err)
	}
	if n.Title != "original title" {
		t.Errorf("expected title to remain %q, got %q", "original title", n.Title)
	}
}

func TestUpdate_AcceptsMultiByteTitleAtMaxRuneLength(t *testing.T) {
	n, err := New("original title", "original", nil, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// "ñ" and "á" are 2 bytes each in UTF-8, so 200 of them is 400 bytes
	// but only 200 runes -- this proves rune-counting, not byte-length, is used.
	title := strings.Repeat("ñá", 100)
	if err := n.Update(title, "updated", nil, time.Now()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.Title != title {
		t.Errorf("expected title to be stored unchanged, got %q", n.Title)
	}
}

func assertTags(t *testing.T, got, want []string) {
	t.Helper()
	if got == nil {
		t.Fatalf("expected non-nil tags %q, got nil", want)
	}
	if !slices.Equal(got, want) {
		t.Errorf("expected tags %q, got %q", want, got)
	}
}

func TestNew_TagNormalization(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"duplicates removed", []string{"go", "go", "backend"}, []string{"go", "backend"}},
		{"trimmed and lowercased", []string{"Go", " go ", "GO"}, []string{"go"}},
		{"empty and whitespace-only discarded", []string{"go", "", "  "}, []string{"go"}},
		{"first-occurrence order preserved", []string{"backend", "go", "Backend", "api"}, []string{"backend", "go", "api"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n, err := New("", "hola", tc.in, time.Now())
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			assertTags(t, n.Tags, tc.want)
		})
	}
}

func TestNew_NilTagsProducesEmptyNonNilSlice(t *testing.T) {
	n, err := New("", "hola", nil, time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assertTags(t, n.Tags, []string{})
}

func TestNew_AcceptsTagAtMaxRuneLength(t *testing.T) {
	tests := []struct {
		name string
		tag  string
	}{
		{"ascii", strings.Repeat("a", 20)},
		// "ñ" and "á" are 2 bytes each: 40 bytes but 20 runes.
		{"multi-byte", strings.Repeat("ñá", 10)},
		// Surrounding whitespace is trimmed before measuring.
		{"padded", "  " + strings.Repeat("a", 20) + "  "},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n, err := New("", "hola", []string{tc.tag}, time.Now())
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			assertTags(t, n.Tags, []string{strings.TrimSpace(tc.tag)})
		})
	}
}

func TestNew_RejectsTagExceedingMaxRuneLength(t *testing.T) {
	_, err := New("", "hola", []string{"go", strings.Repeat("a", 21)}, time.Now())
	if !errors.Is(err, ErrTagTooLong) {
		t.Fatalf("expected ErrTagTooLong, got %v", err)
	}
}

func TestNew_DuplicatesDoNotCountTowardTagLimit(t *testing.T) {
	n, err := New("", "hola", []string{"a", "b", "c", "d", "e", "A"}, time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assertTags(t, n.Tags, []string{"a", "b", "c", "d", "e"})
}

func TestNew_RejectsMoreThanMaxTags(t *testing.T) {
	_, err := New("", "hola", []string{"a", "b", "c", "d", "e", "f"}, time.Now())
	if !errors.Is(err, ErrTooManyTags) {
		t.Fatalf("expected ErrTooManyTags, got %v", err)
	}
}

func TestNew_TagTooLongWinsOverTooManyTags(t *testing.T) {
	// The over-long tag is last, so an implementation that stopped early on
	// the count would return ErrTooManyTags instead.
	tags := []string{"a", "b", "c", "d", "e", "f", strings.Repeat("x", 21)}

	_, err := New("", "hola", tags, time.Now())
	if !errors.Is(err, ErrTagTooLong) {
		t.Fatalf("expected ErrTagTooLong, got %v", err)
	}
}

func TestNew_TagsDoNotAliasInput(t *testing.T) {
	in := []string{"go", "api"}

	n, err := New("", "hola", in, time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	in[0] = "mutated"
	assertTags(t, n.Tags, []string{"go", "api"})
}

func TestUpdate_ReplacesTags(t *testing.T) {
	n, err := New("", "original", []string{"go", "backend"}, time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedAt := time.Now()
	if err := n.Update("", "updated", []string{"API", "api", "db"}, updatedAt); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assertTags(t, n.Tags, []string{"api", "db"})
	if !n.UpdatedAt.Equal(updatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", updatedAt, n.UpdatedAt)
	}
}

func TestUpdate_WithoutTagsClearsThem(t *testing.T) {
	n, err := New("", "original", []string{"go", "backend"}, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := n.Update("", "updated", nil, time.Now()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assertTags(t, n.Tags, []string{})
}

func TestUpdate_RejectedTagsLeaveNoteUnchanged(t *testing.T) {
	tests := []struct {
		name    string
		tags    []string
		wantErr error
	}{
		{"tag too long", []string{strings.Repeat("a", 21)}, ErrTagTooLong},
		{"too many tags", []string{"a", "b", "c", "d", "e", "f"}, ErrTooManyTags},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			created := time.Now().Add(-time.Hour)
			n, err := New("original title", "original", []string{"go", "backend"}, created)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if err := n.Update("new title", "updated", tc.tags, time.Now()); !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected %v, got %v", tc.wantErr, err)
			}
			assertTags(t, n.Tags, []string{"go", "backend"})
			if n.Title != "original title" {
				t.Errorf("expected title to remain %q, got %q", "original title", n.Title)
			}
			if n.Content != "original" {
				t.Errorf("expected content to remain %q, got %q", "original", n.Content)
			}
			if !n.UpdatedAt.Equal(created) {
				t.Errorf("expected UpdatedAt to remain %v, got %v", created, n.UpdatedAt)
			}
		})
	}
}
