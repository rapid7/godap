package filter

import (
   "strings"
)

type BaseFilter struct {
   opts map[string]string
}

func (b *BaseFilter) ParseOpts(args []string) {
   b.opts = make(map[string]string)
   for _, arg := range args {
      params := strings.SplitN(arg, "=", 2)
      if len(params) > 1 {
         b.opts[params[0]] = params[1]
      } else {
         b.opts[params[0]] = ""
      }
   }
}
