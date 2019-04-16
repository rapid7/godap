package filter

import (
	"errors"
	"fmt"
	recog "github.com/hdm/recog-go/pkg/nition"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
	"github.com/rapid7/godap/util"
	"os"
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

// Returns a new recog filter, given the specified `mapped_fields` and database path.
//
// The `mapped_fields` map keys should indicate fields in the source document, and values
// the recog database name to match against.
//
// The `dbpath` can specify an optional directory to recog database files (.xml). If
// not specified, or set to the empty string, a filter will be created based on the content
// embedded in the recog-go package.
func NewFilterRecog(mapped_fields map[string]string, dbpath string) (filterRecog *FilterRecog, err error) {
	filterRecog = new(FilterRecog)

	if mapped_fields == nil {
		return nil, errors.New("Mapped fields must be supplied")
	}
	filterRecog.mapped_fields = mapped_fields

	if dbpath == "" {
		filterRecog.nizer = recog.MustLoadFingerprints()
	} else {
		if filterRecog.nizer, err = recog.LoadFingerprintsDir(dbpath); err != nil {
			filterRecog = nil
		}
	}

	return
}

// Registers this plugin under the "recog" filter name
//
// If the `RECOG_DATABASE_DIRECTORY` environment variable is specified, this plugin
// will use that directory to load recog content. If it is not specified, it will load
// content that is embedded in the `recog-go` package.
func init() {
	factory.RegisterFilter("recog", func(args []string) (filterRecog api.Filter, err error) {
		opts := util.ParseOpts(args)
		filterRecog, err = NewFilterRecog(opts, os.Getenv("RECOG_DATABASE_DIRECTORY"))
		return
	})
}
