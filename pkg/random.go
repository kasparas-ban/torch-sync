package pkg

import (
	"strings"

	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pkg/errors"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	length   = 12
)

// New generates a unique public ID.
func NewRandomID() (string, error) { return nanoid.Generate(alphabet, length) }

// Must is the same as New, but panics on error.
func Must() string { return nanoid.MustGenerate(alphabet, length) }

// Checks if a given field name’s public ID value is valid according to
// the constraints defined by package
func ValidateID(fieldName, id string) error {
	if id == "" {
		return errors.Errorf("%s cannot be blank", fieldName)
	}

	if len(id) != length {
		return errors.Errorf("%s should be %d characters long", fieldName, length)
	}

	if strings.Trim(id, alphabet) != "" {
		return errors.Errorf("%s has invalid characters", fieldName)
	}

	return nil
}
