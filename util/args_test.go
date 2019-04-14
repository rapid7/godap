package util

import (
	"testing"
)

func TestArgsParseLength(t *testing.T) {
	// given
	args := []string{"a", "b=1", "c=abcd"}

	// when
	result := ParseOpts(args)

	// then
	if len(result) != len(args) {
		t.Errorf("Expected %s, received %s", args, result)
	}
}

func TestArgsParseContent(t *testing.T) {
	// given
	args := []string{"a", "b=1", "c=abcd"}

	// when
	result := ParseOpts(args)

	// then
	if result["a"] != "" && result["b"] != "1" && result["c"] != "abcd" {
		t.Errorf("Expected %s, received %v", args, result)
	}
}
