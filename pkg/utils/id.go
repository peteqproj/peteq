package utils

import "github.com/gofrs/uuid"

type (
	// IDGenerator generates ids
	IDGenerator interface {
		GenerateV4() (string, error)
	}

	generator struct{}
)

// NewGenerator builds IDGenerator
func NewGenerator() IDGenerator {
	return &generator{}
}

func (g generator) GenerateV4() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}
