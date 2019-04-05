package util

import (
	"strings"
)

func ParseOpts(args []string) map[string]string {
	parsed_args := make(map[string]string)

	for _, arg := range args {
		params := strings.SplitN(arg, "=", 2)
		if len(params) > 1 {
			parsed_args[params[0]] = params[1]
		} else {
			parsed_args[params[0]] = ""
		}
	}

	return parsed_args
}
