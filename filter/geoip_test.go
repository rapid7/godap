package filter

import (
	"github.com/abh/geoip"
	"github.com/rapid7/godap/factory"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

type GeoIPNoopDecoder struct {
	Decoder
}

func (g *GeoIPNoopDecoder) decode(db *geoip.GeoIP, ip string, field string, doc map[string]interface{}) {
	// noop
}

func TestNewGeoIP(t *testing.T) {
	Convey("Given an empty list of fields", t, func() {
		fields := make([]string, 0)
		decoder := new(GeoIPNoopDecoder)
		geoip_path := "../test/test_data/geoip/GeoIPCity.dat"
		default_database_files := make([]string, 0)

		Convey("When a new geo ip filter is created", func() {
			filter, err := NewFilterGeoIP(fields, geoip_path, default_database_files, decoder)

			Convey("An error is returned indicating the fields were not present", func() {
				So(err, ShouldBeError, ErrNoArgs)
			})

			Convey("The filter returned was nil", func() {
				So(filter, ShouldBeNil)
			})
		})
	})

	Convey("Given a non-existent path to the maxmind database", t, func() {
		fields := []string{"nop"}
		decoder := new(GeoIPNoopDecoder)
		geoip_path := "/some/place/that/doesnt/exist"
		default_database_files := make([]string, 0)

		Convey("When a new geo_ip filter is created", func() {
			filter, err := NewFilterGeoIP(fields, geoip_path, default_database_files, decoder)

			Convey("An error is returned indicating the path doesn't exist", func() {
				So(os.IsNotExist(err), ShouldBeTrue)
			})

			Convey("The filter should be nil", func() {
				So(filter, ShouldBeNil)
			})
		})
	})

	Convey("Given a nil decoder", t, func() {
		fields := []string{"nop"}
		var decoder Decoder = nil
		geoip_path := "../test/test_data/geoip/GeoIPCity.dat"
		default_database_files := make([]string, 0)

		Convey("When a new geo_ip filter is created", func() {
			filter, err := NewFilterGeoIP(fields, geoip_path, default_database_files, decoder)

			Convey("An error is returned indicating the decoder was nil", func() {
				So(err, ShouldBeError)
			})

			Convey("The filter returned was nil", func() {
				So(filter, ShouldBeNil)
			})
		})
	})
}

func TestGeoIPCityDecoder(t *testing.T) {
	Convey("Given an input ip address and geoip database", t, func() {
		ip := "66.92.181.240" // test address in the GeoIPCity.dat test database
		decoder := new(GeoIPCityDecoder)
		database, _ := geoip.Open("../test/test_data/geoip/GeoIPCity.dat")

		Convey("When the IP is processed by the decoder", func() {
			result := make(map[string]interface{})
			decoder.decode(database, ip, "line", result)

			Convey("The result should have expected geo_ip fields", func() {
				So(result, ShouldResemble, map[string]interface{}{
					"line.country_code":   "US",
					"line.country_code3":  "USA",
					"line.country_name":   "United States",
					"line.continent_code": "NA",
					"line.region":         "CA",
					"line.city":           "Fremont",
					"line.postal_code":    "94538",
					"line.latitude":       float32(37.5079),
					"line.longitude":      float32(-121.96),
					"line.dma_code":       807,
					"line.area_code":      510,
				})
			})
		})
	})
}

func TestGeoIPCityFilterFactoryIntegration(t *testing.T) {
	Convey("When a geo_ip (city) filter is created through the factory", t, func() {
		os.Setenv("GEOIP_CITY_DATABASE_PATH", "../test/test_data/geoip/GeoIPCity.dat")
		defer os.Unsetenv("GEOIP_CITY_DATABASE_PATH")
		filterGeoIPCity, err := factory.CreateFilter([]string{"geo_ip", "line"})

		Convey("The geo_ip filter should not be nil", func() {
			So(filterGeoIPCity, ShouldNotBeNil)
		})

		Convey("The error should not be set (nil)", func() {
			So(err, ShouldBeNil)
		})

		Convey("When the input is passed to the filter", func() {
			in_doc := map[string]interface{}{"line": "66.92.181.240"}
			result, err := filterGeoIPCity.Process(in_doc)

			Convey("The result contains only one document", func() {
				So(result, ShouldHaveLength, 1)
			})

			Convey("The first item in the result should have expected geo_ip fields", func() {
				So(result[0], ShouldResemble, map[string]interface{}{
					"line":                "66.92.181.240",
					"line.country_code":   "US",
					"line.country_code3":  "USA",
					"line.country_name":   "United States",
					"line.continent_code": "NA",
					"line.region":         "CA",
					"line.city":           "Fremont",
					"line.postal_code":    "94538",
					"line.latitude":       float32(37.5079),
					"line.longitude":      float32(-121.96),
					"line.dma_code":       807,
					"line.area_code":      510,
				})
			})

			Convey("The error should not be set", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given the GEOIP_CITY_DATABASE_PATH set to an invalid path", t, func() {
		os.Setenv("GEOIP_CITY_DATABASE_PATH", "/does/not/exist")
		defer os.Unsetenv("GEOIP_CITY_DATABASE_PATH")

		Convey("When a geo_ip (city) filter is created", func() {
			filterGeoIPCity, err := factory.CreateFilter([]string{"geo_ip", "line"})

			Convey("The geo_ip filter should be nil", func() {
				So(filterGeoIPCity, ShouldBeNil)
			})

			Convey("The error should indicate the directory does not exist", func() {
				So(os.IsNotExist(err), ShouldBeTrue)
			})
		})
	})
}

func TestGeoIPOrgDecoder(t *testing.T) {
	Convey("Given an input ip address and geoip org database", t, func() {
		ip := "12.87.118.0" // test address in the GeoIPOrg.dat test database
		decoder := new(GeoIPOrgDecoder)
		database, _ := geoip.Open("../test/test_data/geoip/GeoIPOrg.dat")

		Convey("When the IP is processed by the decoder", func() {
			result := make(map[string]interface{})
			decoder.decode(database, ip, "line", result)

			Convey("The result should have expected geo_ip_org fields", func() {
				So(result, ShouldResemble, map[string]interface{}{
					"line.org": "AT&T Worldnet Services",
				})
			})
		})
	})
}

func TestGeoIPOrgFilterFactoryIntegration(t *testing.T) {
	Convey("When a geo_ip_org filter is created through the factory", t, func() {
		os.Setenv("GEOIP_ORG_DATABASE_PATH", "../test/test_data/geoip/GeoIPOrg.dat")
		defer os.Unsetenv("GEOIP_ORG_DATABASE_PATH")
		filterGeoIPOrg, err := factory.CreateFilter([]string{"geo_ip_org", "line"})

		Convey("The geo_ip_org filter should not be nil", func() {
			So(filterGeoIPOrg, ShouldNotBeNil)
		})

		Convey("The error should not be set (nil)", func() {
			So(err, ShouldBeNil)
		})
		Convey("When the input is passed to the filter", func() {
			in_doc := map[string]interface{}{"line": "12.87.118.0"}
			result, err := filterGeoIPOrg.Process(in_doc)

			Convey("The result contains only one document", func() {
				So(result, ShouldHaveLength, 1)
			})

			Convey("The first item in the result should have expected geo_ip_org fields", func() {
				So(result[0], ShouldResemble, map[string]interface{}{
					"line":     "12.87.118.0",
					"line.org": "AT&T Worldnet Services",
				})
			})

			Convey("The error should not be set", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given the GEOIP_ORG_DATABASE_PATH set to an invalid path", t, func() {
		os.Setenv("GEOIP_ORG_DATABASE_PATH", "/does/not/exist")
		defer os.Unsetenv("GEOIP_ORG_DATABASE_PATH")

		Convey("When a geo_ip (city) filter is created", func() {
			filterGeoIPOrg, err := factory.CreateFilter([]string{"geo_ip_org", "line"})

			Convey("The geo_ip filter should be nil", func() {
				So(filterGeoIPOrg, ShouldBeNil)
			})

			Convey("The error should indicate the directory does not exist", func() {
				So(os.IsNotExist(err), ShouldBeTrue)
			})
		})
	})
}
