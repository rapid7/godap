package filter

import (
	"errors"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"github.com/rapid7/godap/api"
	"github.com/rapid7/godap/factory"
	"github.com/rapid7/godap/util"
	"net"
	"os"
)

type GeoIP2CityDecoder struct {
	GeoIPDecoder
	db       *geoip2.Reader
	language string
}

func NewGeoIP2CityDecoder(db *geoip2.Reader, language string) (decoder *GeoIP2CityDecoder, err error) {
	decoder = new(GeoIP2CityDecoder)

	if language == "" {
		return nil, errors.New("A language must be specified.")
	}
	decoder.language = language

	if db == nil {
		return nil, errors.New("A geoip database must be supplied.")
	}
	decoder.db = db
	return
}

func (g *GeoIP2CityDecoder) decode(ip string, field string, doc map[string]interface{}) {
	if parsed_ip := net.ParseIP(ip); parsed_ip != nil {
		if record, err := g.db.City(parsed_ip); err == nil && record != nil {
			doc[fmt.Sprintf("%s.geoip2.city.geoname_id", field)] = record.City.GeoNameID
			doc[fmt.Sprintf("%s.geoip2.city.name", field)] = record.City.Names[g.language]

			doc[fmt.Sprintf("%s.geoip2.continent.code", field)] = record.Continent.Code
			doc[fmt.Sprintf("%s.geoip2.continent.geoname_id", field)] = record.Continent.GeoNameID
			doc[fmt.Sprintf("%s.geoip2.continent.name", field)] = record.Continent.Names[g.language]

			doc[fmt.Sprintf("%s.geoip2.country.geoname_id", field)] = record.Country.GeoNameID
			doc[fmt.Sprintf("%s.geoip2.country.is_eu", field)] = record.Country.IsInEuropeanUnion
			doc[fmt.Sprintf("%s.geoip2.country.iso_code", field)] = record.Country.IsoCode
			doc[fmt.Sprintf("%s.geoip2.country.name", field)] = record.Country.Names[g.language]

			doc[fmt.Sprintf("%s.geoip2.location.accuracy_raidus", field)] = record.Location.AccuracyRadius
			doc[fmt.Sprintf("%s.geoip2.location.latitude", field)] = record.Location.Latitude
			doc[fmt.Sprintf("%s.geoip2.location.longitude", field)] = record.Location.Longitude
			doc[fmt.Sprintf("%s.geoip2.location.metro_code", field)] = record.Location.MetroCode
			doc[fmt.Sprintf("%s.geoip2.location.time_zone", field)] = record.Location.TimeZone

			doc[fmt.Sprintf("%s.geoip2.postal.code", field)] = record.Postal.Code

			doc[fmt.Sprintf("%s.geoip2.registered_country.geoname_id", field)] = record.RegisteredCountry.GeoNameID
			doc[fmt.Sprintf("%s.geoip2.registered_country.is_eu", field)] = record.RegisteredCountry.IsInEuropeanUnion
			doc[fmt.Sprintf("%s.geoip2.registered_country.iso_code", field)] = record.RegisteredCountry.IsoCode
			doc[fmt.Sprintf("%s.geoip2.registered_country.name", field)] = record.RegisteredCountry.Names[g.language]

			doc[fmt.Sprintf("%s.geoip2.represented_country.geoname_id", field)] = record.RepresentedCountry.GeoNameID
			doc[fmt.Sprintf("%s.geoip2.represented_country.is_eu", field)] = record.RepresentedCountry.IsInEuropeanUnion
			doc[fmt.Sprintf("%s.geoip2.represented_country.iso_code", field)] = record.RepresentedCountry.IsoCode
			doc[fmt.Sprintf("%s.geoip2.represented_country.name", field)] = record.RepresentedCountry.Names[g.language]
			doc[fmt.Sprintf("%s.geoip2.represented_country.type", field)] = record.RepresentedCountry.Type

			doc[fmt.Sprintf("%s.geoip2.subdivisions.length", field)] = len(record.Subdivisions)
			for index, subdivision := range record.Subdivisions {
				doc[fmt.Sprintf("%s.geoip2.subdivisions.%d.geoname_id", field, index)] = subdivision.GeoNameID
				doc[fmt.Sprintf("%s.geoip2.subdivisions.%d.iso_code", field, index)] = subdivision.IsoCode
				doc[fmt.Sprintf("%s.geoip2.subdivisions.%d.name", field, index)] = subdivision.Names[g.language]
			}

			doc[fmt.Sprintf("%s.geoip2.traits.is_anon_proxy", field)] = record.Traits.IsAnonymousProxy
			doc[fmt.Sprintf("%s.geoip2.traits.is_satellite", field)] = record.Traits.IsSatelliteProvider
		}
	}
}

type GeoIP2ISPDecoder struct {
	GeoIPDecoder
	db       *geoip2.Reader
	language string
}

func NewGeoIP2ISPDecoder(db *geoip2.Reader, language string) (decoder *GeoIP2ISPDecoder, err error) {
	decoder = new(GeoIP2ISPDecoder)

	if language == "" {
		return nil, errors.New("A language must be specified.")
	}
	decoder.language = language

	if db == nil {
		return nil, errors.New("A geoip database must be supplied.")
	}
	decoder.db = db
	return
}

func (g *GeoIP2ISPDecoder) decode(ip string, field string, doc map[string]interface{}) {
	if parsed_ip := net.ParseIP(ip); parsed_ip != nil {
		if record, err := g.db.ISP(parsed_ip); err == nil && record != nil {
			doc[fmt.Sprintf("%s.geoip2.isp.asn", field)] = record.AutonomousSystemNumber
			doc[fmt.Sprintf("%s.geoip2.isp.asn_org", field)] = record.AutonomousSystemOrganization
			doc[fmt.Sprintf("%s.geoip2.isp.isp", field)] = record.ISP
			doc[fmt.Sprintf("%s.geoip2.isp.org", field)] = record.Organization
		}
	}
}

func open_geoip2_database(geoip_path string, default_database_files []string) (db *geoip2.Reader, err error) {
	stat, err := os.Stat(geoip_path)
	if err != nil {
		return nil, err
	}

	switch mode := stat.Mode(); {
	case mode.IsDir():
		for _, file := range default_database_files {
			db, err = geoip2.Open(fmt.Sprintf("%s/%s", geoip_path, file))
			if err == nil {
				break
			}
		}
	case mode.IsRegular():
		db, err = geoip2.Open(geoip_path)
	}

	return
}

func init() {
	factory.RegisterFilter("geo_ip2_city", func(args []string) (filterGeoIPCity api.Filter, err error) {
		var db *geoip2.Reader
		if db, err = open_geoip2_database(
			util.GetEnv("GEOIP2_CITY_DATABASE_PATH", "/var/lib/geoip2"),
			[]string{"GeoLite2-City.mmdb"}); err != nil {
			return nil, err
		}

		var decoder *GeoIP2CityDecoder
		if decoder, err = NewGeoIP2CityDecoder(db, util.GetEnv("GEOIP2_LANGUAGE", "en")); err != nil {
			return nil, err
		}

		return NewFilterGeoIP(args, decoder)
	})
}

func init() {
	factory.RegisterFilter("geo_ip2_isp", func(args []string) (filterGeoIPCity api.Filter, err error) {
		var db *geoip2.Reader
		if db, err = open_geoip2_database(
			util.GetEnv("GEOIP2_ISP_DATABASE_PATH", "/var/lib/geoip2"),
			[]string{"GeoLite2-ISP.mmdb"}); err != nil {
			return nil, err
		}

		var decoder *GeoIP2ISPDecoder
		if decoder, err = NewGeoIP2ISPDecoder(db, util.GetEnv("GEOIP2_LANGUAGE", "en")); err != nil {
			return nil, err
		}

		return NewFilterGeoIP(args, decoder)
	})
}
