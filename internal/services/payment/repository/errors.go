package repository

import "errors"

var (
	ErrUuidInvalidFormat = errors.New("uuid has invalid format")
	ErrNotFound          = errors.New("record is not found")
)
