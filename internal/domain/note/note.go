package note

import "time"

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
	n.Title = title
	n.Content = content
	n.UpdatedAt = now
	return nil
}
