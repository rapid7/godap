package filter

import (
	"strings"
	"testing"
)

// Ensures that the recog library uses the file name of the recog
// fingerprint database. This condition is tested by looking for
// _any_ recog field within the result document returned by the
// filter.
func TestUsesFilenameForDatabaseName(t *testing.T) {
	// given
	filterRecog := NewFilterRecog(map[string]string{
		"input": "dns_versionbind.xml",
	})

	// when
	result, _ := filterRecog.Process(map[string]interface{}{
		"input": "9.8.2rc1-RedHat-9.8.2-0.62.rc1.el6_9.2",
	})

	// then
	has_recog := false
	for k, _ := range result[0] {
		if strings.Contains(k, ".recog.") {
			has_recog = true
		}
	}
	if !has_recog {
		t.Errorf("No recog fields present in %s", result)
	}
}
