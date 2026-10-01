package note

import (
	"time"
	"unicode/utf8"
)

const maxTitleRunes = 200

type Note struct {
	ID        string
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(title, content string, now time.Time) (*Note, error) {
	if content == "" {
		return nil, ErrEmptyContent
	}
	if utf8.RuneCountInString(title) > maxTitleRunes {
		return nil, ErrTitleTooLong
	}
	return &Note{
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (n *Note) Update(title, content string, now time.Time) error {
	if content == "" {
		return ErrEmptyContent
	}
	if utf8.RuneCountInString(title) > maxTitleRunes {
		return ErrTitleTooLong
	}
	n.Title = title
	n.Content = content
	n.UpdatedAt = now
	return nil
}
