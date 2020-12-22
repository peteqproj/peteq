package errors

import "fmt"

type (
	// NotFoundError common error when something was not found
	NotFoundError struct {
		Resource string
		ID       string
	}
)

// NewNotFoundError builds NotFoundError
func NewNotFoundError(resource string, id string) NotFoundError {
	return NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

func (n NotFoundError) Error() string {
	if n.ID != "" {
		return fmt.Sprintf("Resource %s with ID %s was not found", n.Resource, n.ID)
	}
	return fmt.Sprintf("Resource %s was not found", n.Resource)
}
