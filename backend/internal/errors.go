package internal

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("record not found")

type NotFoundError struct {
	Table   string
	Key     string
	Value   string
	Message string
}

func (e NotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Key != "" && e.Value != "" {
		return fmt.Sprintf("unable to find %s with %s=%s", e.Table, e.Key, e.Value)
	}
	return ErrNotFound.Error()
}

func (e NotFoundError) Is(target error) bool {
	return errors.Is(target, ErrNotFound)
}

func NewNotFoundError(table, key, value, message string) NotFoundError {
	return NotFoundError{
		Table:   table,
		Key:     key,
		Value:   value,
		Message: message,
	}
}

func EarlyApplicationFailed(title, action string) string {
	result := `
-----------------------------------------
Application Failed to Start
-----------------------------------------

# What's wrong?
%s

# How to fix it?
%s
`

	result = fmt.Sprintf(result, title, action)
	return result
}
