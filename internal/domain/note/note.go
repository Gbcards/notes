package note

import "time"

type Note struct {
	ID        string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(content string, now time.Time) (*Note, error) {
	if content == "" {
		return nil, ErrEmptyContent
	}
	return &Note{
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (n *Note) UpdateContent(content string, now time.Time) error {
	if content == "" {
		return ErrEmptyContent
	}
	n.Content = content
	n.UpdatedAt = now
	return nil
}
