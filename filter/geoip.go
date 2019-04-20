package filter

import (
	"errors"
	"github.com/rapid7/godap/api"
)

type GeoIPDecoder interface {
	decode(ip string, field string, doc map[string]interface{})
}

type FilterGeoIP struct {
	api.Filter
	decoder GeoIPDecoder
	fields  []string
}

func (filterGeoIP *FilterGeoIP) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
	for _, k := range filterGeoIP.fields {
		if docv, ok := doc[k]; ok {
			filterGeoIP.decoder.decode(docv.(string), k, doc)
		}
	}
	return []map[string]interface{}{doc}, nil
}

func NewFilterGeoIP(fields []string, decoder GeoIPDecoder) (filterGeoIP *FilterGeoIP, err error) {
	filterGeoIP = new(FilterGeoIP)

	if decoder == nil {
		return nil, errors.New("Nil Geo IP decoder specified")
	}
	filterGeoIP.decoder = decoder

	if fields == nil || len(fields) == 0 {
		return nil, ErrNoArgs
	}
	filterGeoIP.fields = fields

	return
}
