package errors

import ()

// Error interface duplicates the standard library errors
type Error interface {
	Error() string
}

func New(msg string) Error {
	return StandardError{msg}
}
