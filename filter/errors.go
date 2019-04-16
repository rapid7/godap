package filter

import (
	"errors"
)

var (
	ErrNoArgs         = errors.New("Arguments must be supplied")
	ErrInvalidMapArgs = errors.New("Expected two parameters in the format <key>=<value>")
)
