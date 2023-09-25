package payment

import "errors"

var (
	ErrUuidInvalidFormat = errors.New("uuid has invalid format")
	ErrNotFound          = errors.New("record not found")
	ErrUnexpectedResult  = errors.New("unexpected error")
	ErrProvider          = errors.New("something happened on the provider side")
	ErrStore             = errors.New("something happened on the stores side")
)
