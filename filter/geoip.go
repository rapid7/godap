// +build libgeoip

package filter

import (
	"errors"
	"fmt"
	"github.com/abh/geoip"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
	"github.com/rapid7/godap/util"
	"os"
)

type GeoIPDecoder interface {
	decode(db *geoip.GeoIP, ip string, field string, doc map[string]interface{})
}

type FilterGeoIP struct {
	api.Filter
	decoder GeoIPDecoder
	db      *geoip.GeoIP
	fields  []string
}

func (filterGeoIP *FilterGeoIP) Process(doc map[string]interface{}) (res []map[string]interface{}, err error) {
	for _, k := range filterGeoIP.fields {
		if docv, ok := doc[k]; ok {
			filterGeoIP.decoder.decode(filterGeoIP.db, docv.(string), k, doc)
		}
	}
	return []map[string]interface{}{doc}, nil
}

func NewFilterGeoIP(fields []string, geoip_path string, default_database_files []string, decoder GeoIPDecoder) (filterGeoIP *FilterGeoIP, err error) {
	filterGeoIP = new(FilterGeoIP)
	filterGeoIP.fields = fields

	if decoder == nil {
		return nil, errors.New("Nil Geo IP decoder specified")
	}
	filterGeoIP.decoder = decoder

	if len(fields) == 0 {
		return nil, ErrNoArgs
	}

	stat, err := os.Stat(geoip_path)
	if err != nil {
		return nil, err
	}

	switch mode := stat.Mode(); {
	case mode.IsDir():
		for _, file := range default_database_files {
			filterGeoIP.db, err = geoip.Open(fmt.Sprintf("%s/%s", geoip_path, file))
			if err == nil {
				break
			}
		}
	case mode.IsRegular():
		filterGeoIP.db, err = geoip.Open(geoip_path)
	}

	return
}

type GeoIPCityDecoder struct {
	GeoIPDecoder
}

func (g *GeoIPCityDecoder) decode(db *geoip.GeoIP, ip string, field string, doc map[string]interface{}) {
	record := db.GetRecord(ip)
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
	factory.RegisterFilter("geo_ip", func(args []string) (filterGeoIPCity api.Filter, err error) {
		filterGeoIPCity, err = NewFilterGeoIP(
			args,
			util.GetEnv("GEOIP_CITY_DATABASE_PATH", "/var/lib/geoip"),
			[]string{"geoip.dat", "geoip_city.dat", "GeoCity.dat", "IP_V4_CITY.dat", "GeoCityLite.dat"},
			new(GeoIPCityDecoder))
		return
	})
}

type GeoIPOrgDecoder struct {
	GeoIPDecoder
}

func (g *GeoIPOrgDecoder) decode(db *geoip.GeoIP, ip string, field string, doc map[string]interface{}) {
	record := db.GetOrg(ip)
	if record != "" {
		doc[fmt.Sprintf("%s.org", field)] = record
	}
}

func init() {
	factory.RegisterFilter("geo_ip_org", func(args []string) (filterGeoIPOrg api.Filter, err error) {
		filterGeoIPOrg, err = NewFilterGeoIP(
			args,
			util.GetEnv("GEOIP_ORG_DATABASE_PATH", "/var/lib/geoip"),
			[]string{"geoip_org.dat", "IP_V4_ORG.dat"},
			new(GeoIPOrgDecoder))
		return
	})
}
