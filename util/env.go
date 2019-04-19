package util

import (
	"os"
)

// Returns the content of the environment variable named by `key`, falling back to
// `fallback` if the environment variable does not exist.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
