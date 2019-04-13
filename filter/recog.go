package filter

import (
	"fmt"
	recog "github.com/hdm/recog-go/pkg/nition"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
	"github.com/rapid7/godap/util"
)

type FilterRecog struct {
	api.Filter
	nizer         *recog.FingerprintSet
	mapped_fields map[string]string
}

func (fr *FilterRecog) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
	for match_field, match_db := range fr.mapped_fields {
		if field_value_untyped, ok := doc[match_field]; ok {
			if field_value, ok := field_value_untyped.(string); ok {
				m := fr.nizer.MatchFirst(match_db, field_value)
				if m.Matched {
					for k, v := range m.Values {
						doc[fmt.Sprintf("%s.recog.%s", match_field, k)] = v
					}
				}
			}

		}
	}
	return []map[string]interface{}{doc}, nil
}

func NewFilterRecog(mapped_fields map[string]string) *FilterRecog {
	filterRecog := new(FilterRecog)
	filterRecog.mapped_fields = mapped_fields
	filterRecog.nizer = recog.MustLoadFingerprints()
	return filterRecog
}

func init() {
	factory.RegisterFilter("recog", func(args []string) (lines api.Filter, err error) {
		opts := util.ParseOpts(args)
		filterRecog := NewFilterRecog(opts)
		return filterRecog, nil
	})
}
