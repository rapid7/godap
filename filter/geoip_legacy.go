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

func open_geoip_database(geoip_path string, default_database_files []string) (db *geoip.GeoIP, err error) {
	stat, err := os.Stat(geoip_path)
	if err != nil {
		return nil, err
	}

	switch mode := stat.Mode(); {
	case mode.IsDir():
		for _, file := range default_database_files {
			db, err = geoip.Open(fmt.Sprintf("%s/%s", geoip_path, file))
			if err == nil {
				break
			}
		}
	case mode.IsRegular():
		db, err = geoip.Open(geoip_path)
	}

	return
}

type GeoIPCityDecoder struct {
	GeoIPDecoder
	db *geoip.GeoIP
}

func NewGeoIPCityDecoder(db *geoip.GeoIP) (decoder *GeoIPCityDecoder, err error) {
	decoder = new(GeoIPCityDecoder)

	if db == nil {
		return nil, errors.New("A geoip database must be supplied.")
	}
	decoder.db = db
	return
}

func (g *GeoIPCityDecoder) decode(ip string, field string, doc map[string]interface{}) {
	record := g.db.GetRecord(ip)
	if record != nil {
		doc[fmt.Sprintf("%s.country_code", field)] = record.CountryCode
		doc[fmt.Sprintf("%s.country_code3", field)] = record.CountryCode3
		doc[fmt.Sprintf("%s.country_name", field)] = record.CountryName
		doc[fmt.Sprintf("%s.region", field)] = record.Region
		doc[fmt.Sprintf("%s.region_name", field)] = geoip.GetRegionName(record.Region, record.CountryCode)
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
		var db *geoip.GeoIP
		if db, err = open_geoip_database(
			util.GetEnv("GEOIP_CITY_DATABASE_PATH", "/var/lib/geoip"),
			[]string{"geoip.dat", "geoip_city.dat", "GeoCity.dat", "IP_V4_CITY.dat", "GeoCityLite.dat"}); err != nil {
			return nil, err
		}

		var decoder GeoIPDecoder
		if decoder, err = NewGeoIPCityDecoder(db); err != nil {
			return nil, err
		}
		return NewFilterGeoIP(args, decoder)
	})
}

type GeoIPOrgDecoder struct {
	GeoIPDecoder
	db *geoip.GeoIP
}

func NewGeoIPOrgDecoder(db *geoip.GeoIP) (decoder *GeoIPOrgDecoder, err error) {
	decoder = new(GeoIPOrgDecoder)

	if db == nil {
		return nil, errors.New("A geoip database must be supplied.")
	}
	decoder.db = db
	return
}

func (g *GeoIPOrgDecoder) decode(ip string, field string, doc map[string]interface{}) {
	record := g.db.GetOrg(ip)
	if record != "" {
		doc[fmt.Sprintf("%s.org", field)] = record
	}
}

func init() {
	factory.RegisterFilter("geo_ip_org", func(args []string) (filterGeoIPOrg api.Filter, err error) {
		var db *geoip.GeoIP
		if db, err = open_geoip_database(
			util.GetEnv("GEOIP_ORG_DATABASE_PATH", "/var/lib/geoip"),
			[]string{"geoip_org.dat", "IP_V4_ORG.dat"}); err != nil {
			return nil, err
		}

		var decoder GeoIPDecoder
		if decoder, err = NewGeoIPOrgDecoder(db); err != nil {
			return nil, err
		}

		return NewFilterGeoIP(args, decoder)
	})
}
