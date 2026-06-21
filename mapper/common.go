package mapper

import (
	"fmt"

	"github.com/google/uuid"
)

// BoolPtrToBool converts a bool pointer to bool; returns false if the pointer is nil.
func BoolPtrToBool(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}

// ParseUUIDParam parses a UUID from an optional query or path parameter string.
// A nil pointer is treated as uuid.Nil with no error, preserving the existing
// nil-guard semantics at call sites. A non-nil but malformed value returns an
// error instead of panicking (replacing uuid.MustParse on user input).
func ParseUUIDParam(s *string) (uuid.UUID, error) {
	if s == nil {
		return uuid.Nil, nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID parameter %q: %w", *s, err)
	}
	return id, nil
}
