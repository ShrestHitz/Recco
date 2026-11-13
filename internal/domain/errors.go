package domain

import "errors"

var (
	// ErrNotFound is returned when a requested entity does not exist.
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput is returned when input validation fails.
	ErrInvalidInput = errors.New("invalid input")

	// ErrConflict is returned when an entity violates uniqueness constraints.
	ErrConflict = errors.New("resource conflict or already exists")

	// ErrUnauthorized is returned when authentication/authorization fails.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInvalidStatusTransition is returned when an illegal state change is requested.
	ErrInvalidStatusTransition = errors.New("invalid shelf status transition")
)
