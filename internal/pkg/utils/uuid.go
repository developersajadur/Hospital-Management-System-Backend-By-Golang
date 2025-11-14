package utils

import "github.com/google/uuid"

// UUIDPtr converts *string → *uuid.UUID safely.
// Returns nil if input is nil, empty, or invalid.
func UUIDPtr(id *string) *uuid.UUID {
	if id == nil || *id == "" {
		return nil
	}

	parsed, err := uuid.Parse(*id)
	if err != nil {
		return nil
	}

	return &parsed
}
