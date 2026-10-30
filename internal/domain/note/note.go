package note

import (
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxTitleRunes = 200
	maxTagRunes   = 20
	maxTags       = 5
)

type Note struct {
	ID        string
	Title     string
	Content   string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(title, content string, tags []string, now time.Time) (*Note, error) {
	if content == "" {
		return nil, ErrEmptyContent
	}
	if utf8.RuneCountInString(title) > maxTitleRunes {
		return nil, ErrTitleTooLong
	}
	normalized, err := normalizeTags(tags)
	if err != nil {
		return nil, err
	}
	return &Note{
		Title:     title,
		Content:   content,
		Tags:      normalized,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (n *Note) Update(title, content string, tags []string, now time.Time) error {
	if content == "" {
		return ErrEmptyContent
	}
	if utf8.RuneCountInString(title) > maxTitleRunes {
		return ErrTitleTooLong
	}
	normalized, err := normalizeTags(tags)
	if err != nil {
		return err
	}
	n.Title = title
	n.Content = content
	n.Tags = normalized
	n.UpdatedAt = now
	return nil
}

// normalizeTags trims and lowercases each tag, drops empty entries and
// duplicates (keeping first-occurrence order), and enforces the per-tag
// length and total count limits. The count is checked after deduplication,
// and the loop never stops early on the count, so ErrTagTooLong takes
// precedence over ErrTooManyTags. The result is always a newly allocated,
// non-nil slice that never aliases raw.
func normalizeTags(raw []string) ([]string, error) {
	result := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, t := range raw {
		tag := strings.ToLower(strings.TrimSpace(t))
		if tag == "" {
			continue
		}
		if utf8.RuneCountInString(tag) > maxTagRunes {
			return nil, ErrTagTooLong
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		result = append(result, tag)
	}
	if len(result) > maxTags {
		return nil, ErrTooManyTags
	}
	return result, nil
}
