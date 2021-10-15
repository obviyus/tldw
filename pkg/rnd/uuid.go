package rnd

import (
	uuid "github.com/satori/go.uuid"
)

// UUID returns a standard, random UUID as string.
func UUID() string {
	return uuid.NewV4().String()
}

// IsUUID Returns true if the string looks like a standard UUID.
func IsUUID(s string) bool {
	return len(s) == 36 && IsHex(s)
}
