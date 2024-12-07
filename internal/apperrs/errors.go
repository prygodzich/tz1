package apperrs

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrNotInitialized = errors.New("not initialized")
)
