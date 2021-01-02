package db

import (
	"strings"
)

const acceptedChars = "abcdefghijklmnopqrstuvwxyz_0123456789"

// SchemaIsSafe tests the safety of the schema
func SchemaIsSafe(schema string) bool {
	for _, c := range schema {
		if !strings.Contains(acceptedChars, string(c)) {
			return false
		}
	}
	return true
}
