package filter

import (
   "testing"
)

func TestParseKey(t *testing.T) {
   baseFilter := new(BaseFilter)
   baseFilter.ParseOpts([]string{"foo"})
   if baseFilter.opts["foo"] != "" {
      t.Errorf("Expected the empty string, received %s", baseFilter.opts["foo"])
   }
}

func TestParseKeyValue(t *testing.T) {
   baseFilter := new(BaseFilter)
   baseFilter.ParseOpts([]string{"foo=bar"})
   if baseFilter.opts["foo"] != "bar" {
      t.Errorf("Expected the value \"bar\", received %s", baseFilter.opts["foo"])
   }
}
