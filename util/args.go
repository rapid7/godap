package util

import (
	"strings"
)

// Returns a map of options parsed from an input string array.
//
// The input string array may be formatted as "a=b" (key/value) or "a" (no value).
// This method will then return a map containing the parsed arguments separated
// by an equals sign, with the right hand side content being the value for the
// left hand side (key).
func ParseOpts(args []string) map[string]string {
	opts := make(map[string]string)
	for _, arg := range args {
		params := strings.SplitN(arg, "=", 2)
		if len(params) > 1 {
			opts[params[0]] = params[1]
		} else {
			opts[params[0]] = ""
		}
	}
	return opts
}
