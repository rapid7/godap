package filter

import (
  "fmt"
  "github.com/abh/geoip"
  "github.com/rapid7/godap/api"
  "github.com/rapid7/godap/factory"
)

type Decoder interface {
  Decode(string, string, map[string]interface{})
}

/////////////////////////////////////////////////
// geo_ip filter
/////////////////////////////////////////////////

type FilterGeoIP struct {
  db *geoip.GeoIP
  BaseFilter
  Decoder
}

func (b *FilterGeoIP) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
  for k, _ := range b.opts {
    if docv, ok := doc[k]; ok {
      b.Decode(docv.(string), k, doc)
    }
  }
  return []map[string]interface{}{doc}, nil
}

func (f *FilterGeoIP) Decode(ip string, field string, doc map[string]interface{}) {
  record := f.db.GetRecord(ip)
  if record != nil {
    doc[fmt.Sprintf("%s.country_code", field)] = record.CountryCode
    doc[fmt.Sprintf("%s.country_code3", field)] = record.CountryCode3
    doc[fmt.Sprintf("%s.country_name", field)] = record.CountryName
    doc[fmt.Sprintf("%s.continent_code", field)] = record.ContinentCode
    doc[fmt.Sprintf("%s.region", field)] = record.Region
    doc[fmt.Sprintf("%s.city", field)] = record.City
    doc[fmt.Sprintf("%s.postal_code", field)] = record.PostalCode
    doc[fmt.Sprintf("%s.latitude", field)] = record.Latitude
    doc[fmt.Sprintf("%s.longitude", field)] = record.Longitude
    doc[fmt.Sprintf("%s.dma_code", field)] = record.MetroCode
    doc[fmt.Sprintf("%s.area_code", field)] = record.AreaCode
  }
}

func init() {
  factory.RegisterFilter("geo_ip", func(args []string) (lines api.Filter, err error) {
    filterGeoIP := &FilterGeoIP{}
    filterGeoIP.ParseOpts(args)
    for _, file := range []string{"geoip.dat", "geoip_city.dat", "GeoCity.dat", "IP_V4_CITY.dat", "GeoCityLite.dat"} {
      filterGeoIP.db, err = geoip.Open(fmt.Sprintf("%s/%s", "/var/lib/geoip", file))
      if err == nil {
        break
      }
    }
    if filterGeoIP == nil {
      err = fmt.Errorf("Could not open geoip database")
      return nil, err
    }
    return filterGeoIP, nil
  })
}
