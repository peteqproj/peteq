package errors

import "errors"

// New warps errors.New
func New(text string) error {
	return errors.New(text)
}

// As warps errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Is warps errors.As
func Is(err, target error) bool {
	return errors.Is(err, target)
}
