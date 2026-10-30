package note

import "errors"

var (
	ErrNotFound     = errors.New("note not found")
	ErrEmptyContent = errors.New("note content must not be empty")
	ErrTitleTooLong = errors.New("note title must not exceed 200 characters")
	ErrTagTooLong   = errors.New("note tag must not exceed 20 characters")
	ErrTooManyTags  = errors.New("note must not have more than 5 tags")
)
