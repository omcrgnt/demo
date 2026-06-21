package domain

import "errors"

var (
	ErrNotFound     = errors.New("item not found")
	ErrInvalidInput = errors.New("invalid input")
)
