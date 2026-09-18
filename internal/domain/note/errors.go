package note

import "errors"

var (
	ErrNotFound     = errors.New("note not found")
	ErrEmptyContent = errors.New("note content must not be empty")
)
