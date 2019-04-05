package filter

import (
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/util"
)

type BaseFilter struct {
	opts map[string]string
	api.Filter
}

func (b *BaseFilter) ParseOpts(args []string) {
	b.opts = util.ParseOpts(args)
}
