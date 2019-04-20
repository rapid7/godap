package filter

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestGeoIP2Database(t *testing.T) {
	Convey("Given a non-existent path to the maxmind database", t, func() {
		geoip_path := "/some/place/that/doesnt/exist"
		default_database_files := make([]string, 0)

		Convey("When the database is opened", func() {
			filter, err := open_geoip2_database(geoip_path, default_database_files)

			Convey("An error is returned indicating the path doesn't exist", func() {
				So(os.IsNotExist(err), ShouldBeTrue)
			})

			Convey("The filter should be nil", func() {
				So(filter, ShouldBeNil)
			})
		})
	})

	Convey("Given a path to an existing maxmind database", t, func() {
		geoip_path := "../test/test_data/geoip2/GeoIP2-City-Test.mmdb"
		default_database_files := make([]string, 0)

		Convey("When the database is opened", func() {
			filter, err := open_geoip2_database(geoip_path, default_database_files)

			Convey("No error returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("The filter should not be nil", func() {
				So(filter, ShouldNotBeNil)
			})
		})
	})
}

func TestGeoIP2CityDecoder(t *testing.T) {
	Convey("Given an input ip address and geoip database", t, func() {
		ip := "81.2.69.142" // test address in the GeoIP2-City-Test.mmdb test database
		db, _ := open_geoip2_database("../test/test_data/geoip2/GeoIP2-City-Test.mmdb", []string{})
		decoder, _ := NewGeoIP2CityDecoder(db, "en")

		Convey("When the IP is processed by the decoder", func() {
			result := make(map[string]interface{})
			decoder.decode(ip, "line", result)

			Convey("The result should have expected geo_ip2 city fields", func() {
				So(result, ShouldHaveLength, 30)
				So(result["line.geoip2.city.geoname_id"], ShouldEqual, 2643743)
				So(result["line.geoip2.city.name"], ShouldEqual, "London")
				So(result["line.geoip2.continent.code"], ShouldEqual, "EU")
				So(result["line.geoip2.continent.geoname_id"], ShouldEqual, 6255148)
				So(result["line.geoip2.continent.name"], ShouldEqual, "Europe")
				So(result["line.geoip2.country.geoname_id"], ShouldEqual, 2635167)
				So(result["line.geoip2.country.is_eu"], ShouldBeTrue)
				So(result["line.geoip2.country.iso_code"], ShouldEqual, "GB")
				So(result["line.geoip2.country.name"], ShouldEqual, "United Kingdom")
				So(result["line.geoip2.location.accuracy_raidus"], ShouldEqual, 10)
				So(result["line.geoip2.location.latitude"], ShouldAlmostEqual, 51.5142)
				So(result["line.geoip2.location.longitude"], ShouldAlmostEqual, -0.0931)
				So(result["line.geoip2.location.metro_code"], ShouldEqual, 0)
				So(result["line.geoip2.location.time_zone"], ShouldEqual, "Europe/London")
				So(result["line.geoip2.postal.code"], ShouldBeEmpty)
				So(result["line.geoip2.registered_country.geoname_id"], ShouldEqual, 6252001)
				So(result["line.geoip2.registered_country.is_eu"], ShouldBeFalse)
				So(result["line.geoip2.registered_country.iso_code"], ShouldEqual, "US")
				So(result["line.geoip2.registered_country.name"], ShouldEqual, "United States")
				So(result["line.geoip2.represented_country.geoname_id"], ShouldEqual, 0)
				So(result["line.geoip2.represented_country.is_eu"], ShouldBeFalse)
				So(result["line.geoip2.represented_country.iso_code"], ShouldBeEmpty)
				So(result["line.geoip2.represented_country.name"], ShouldBeEmpty)
				So(result["line.geoip2.represented_country.type"], ShouldBeEmpty)
				So(result["line.geoip2.subdivisions.length"], ShouldEqual, 1)
				So(result["line.geoip2.subdivisions.0.geoname_id"], ShouldEqual, 6269131)
				So(result["line.geoip2.subdivisions.0.iso_code"], ShouldEqual, "ENG")
				So(result["line.geoip2.subdivisions.0.name"], ShouldEqual, "England")
				So(result["line.geoip2.traits.is_anon_proxy"], ShouldBeFalse)
				So(result["line.geoip2.traits.is_satellite"], ShouldBeFalse)
			})
		})
	})
}

func TestGeoIP2ISPDecoder(t *testing.T) {
	Convey("Given an input ip address and geoip database", t, func() {
		ip := "1.128.0.0" // test address in the GeoIP2-City-Test.mmdb test database
		db, _ := open_geoip2_database("../test/test_data/geoip2/GeoIP2-ISP-Test.mmdb", []string{})
		decoder, _ := NewGeoIP2ISPDecoder(db, "en")

		Convey("When the IP is processed by the decoder", func() {
			result := make(map[string]interface{})
			decoder.decode(ip, "line", result)

			Convey("The result should have expected geo_ip2 isp fields", func() {
				So(result, ShouldHaveLength, 4)
				So(result["line.geoip2.isp.asn"], ShouldEqual, 1221)
				So(result["line.geoip2.isp.asn_org"], ShouldEqual, "Telstra Pty Ltd")
				So(result["line.geoip2.isp.isp"], ShouldEqual, "Telstra Internet")
				So(result["line.geoip2.isp.org"], ShouldEqual, "Telstra Internet")
			})
		})
	})
}
