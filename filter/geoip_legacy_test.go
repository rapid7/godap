package filter

import (
	"github.com/rapid7/godap/factory"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestGeoIPLegacyDatabase(t *testing.T) {
	Convey("Given a non-existent path to the maxmind database", t, func() {
		geoip_path := "/some/place/that/doesnt/exist"
		default_database_files := make([]string, 0)

		Convey("When the database is opened", func() {
			filter, err := open_geoip_database(geoip_path, default_database_files)

			Convey("An error is returned indicating the path doesn't exist", func() {
				So(os.IsNotExist(err), ShouldBeTrue)
			})

			Convey("The filter should be nil", func() {
				So(filter, ShouldBeNil)
			})
		})
	})
}

func TestGeoIPCityDecoder(t *testing.T) {
	Convey("Given an input ip address and geoip database", t, func() {
		ip := "66.92.181.240" // test address in the GeoIPCity.dat test database
		db, _ := open_geoip_database("../test/test_data/geoip/GeoIPCity.dat", []string{})
		decoder, _ := NewGeoIPCityDecoder(db)

		Convey("When the IP is processed by the decoder", func() {
			result := make(map[string]interface{})
			decoder.decode(ip, "line", result)

			Convey("The result should have expected geo_ip fields", func() {
				So(result, ShouldResemble, map[string]interface{}{
					"line.country_code":  "US",
					"line.country_code3": "USA",
					"line.country_name":  "United States",
					"line.region":        "CA",
					"line.region_name":   "California",
					"line.city":          "Fremont",
					"line.postal_code":   "94538",
					"line.latitude":      "37.50790023803711",
					"line.longitude":     "-121.95999908447266",
					"line.dma_code":      "807",
					"line.area_code":     "510",
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
					"line":               "66.92.181.240",
					"line.country_code":  "US",
					"line.country_code3": "USA",
					"line.country_name":  "United States",
					"line.region":        "CA",
					"line.region_name":   "California",
					"line.city":          "Fremont",
					"line.postal_code":   "94538",
					"line.latitude":      "37.50790023803711",
					"line.longitude":     "-121.95999908447266",
					"line.dma_code":      "807",
					"line.area_code":     "510",
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
		db, _ := open_geoip_database("../test/test_data/geoip/GeoIPOrg.dat", []string{})
		decoder, _ := NewGeoIPOrgDecoder(db)

		Convey("When the IP is processed by the decoder", func() {
			result := make(map[string]interface{})
			decoder.decode(ip, "line", result)

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

		Convey("When a geo_ip_org filter is created", func() {
			filterGeoIPOrg, err := factory.CreateFilter([]string{"geo_ip_org", "line"})

			Convey("The geo_ip_org filter should be nil", func() {
				So(filterGeoIPOrg, ShouldBeNil)
			})

			Convey("The error should indicate the directory does not exist", func() {
				So(os.IsNotExist(err), ShouldBeTrue)
			})
		})
	})
}

func TestGeoIPASNDecoder(t *testing.T) {
	Convey("Given an input ip address and geoip asn database", t, func() {
		ip := "1.128.0.0" // test address in the GeoIPOrg.dat test database
		db, _ := open_geoip_database("../test/test_data/geoip/GeoIPASNum.dat", []string{})
		decoder, _ := NewGeoIPASNDecoder(db)

		Convey("When the IP is processed by the decoder", func() {
			result := make(map[string]interface{})
			decoder.decode(ip, "line", result)

			Convey("The result should have the expected `asn` field", func() {
				So(result, ShouldResemble, map[string]interface{}{
					"line.asn": "AS1221",
				})
			})
		})
	})
}

func TestGeoIPASNFilterFactoryIntegration(t *testing.T) {
	Convey("When a geo_ip_asn filter is created through the factory", t, func() {
		os.Setenv("GEOIP_ASN_DATABASE_PATH", "../test/test_data/geoip/GeoIPASNum.dat")
		defer os.Unsetenv("GEOIP_ASN_DATABASE_PATH")
		filterGeoIPASN, err := factory.CreateFilter([]string{"geo_ip_asn", "line"})

		Convey("The geo_ip_asn filter should not be nil", func() {
			So(filterGeoIPASN, ShouldNotBeNil)
		})

		Convey("The error should not be set (nil)", func() {
			So(err, ShouldBeNil)
		})
		Convey("When the input is passed to the filter", func() {
			in_doc := map[string]interface{}{"line": "1.128.0.0"}
			result, err := filterGeoIPASN.Process(in_doc)

			Convey("The result contains only one document", func() {
				So(result, ShouldHaveLength, 1)
			})

			Convey("The first item in the result should have expected geo_ip_org fields", func() {
				So(result[0], ShouldResemble, map[string]interface{}{
					"line":     "1.128.0.0",
					"line.asn": "AS1221",
				})
			})

			Convey("The error should not be set", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given the GEOIP_ASN_DATABASE_PATH set to an invalid path", t, func() {
		os.Setenv("GEOIP_ASN_DATABASE_PATH", "/does/not/exist")
		defer os.Unsetenv("GEOIP_ASN_DATABASE_PATH")

		Convey("When a geo_ip_asn filter is created", func() {
			filterGeoIPOrg, err := factory.CreateFilter([]string{"geo_ip_org", "line"})

			Convey("The geo_ip_asn filter should be nil", func() {
				So(filterGeoIPOrg, ShouldBeNil)
			})

			Convey("The error should indicate the directory does not exist", func() {
				So(os.IsNotExist(err), ShouldBeTrue)
			})
		})
	})
}
