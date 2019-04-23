package filter

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type GeoIPNoopDecoder struct {
	GeoIPDecoder
}

func (g *GeoIPNoopDecoder) decode(ip string, field string, doc map[string]interface{}) {
	// noop
}

func TestNewGeoIP(t *testing.T) {
	Convey("Given a nil list of fields", t, func() {
		var fields []string = nil
		decoder := new(GeoIPNoopDecoder)

		Convey("When a new geo ip filter is created", func() {
			filter, err := NewFilterGeoIP(fields, decoder)

			Convey("An error is returned indicating the fields parameter was nil", func() {
				So(err, ShouldBeError)
			})

			Convey("The filter returned was nil", func() {
				So(filter, ShouldBeNil)
			})
		})
	})

	Convey("Given an empty list of fields", t, func() {
		fields := make([]string, 0)
		decoder := new(GeoIPNoopDecoder)

		Convey("When a new geo ip filter is created", func() {
			filter, err := NewFilterGeoIP(fields, decoder)

			Convey("An error is returned indicating the fields were not present", func() {
				So(err, ShouldBeError, ErrNoArgs)
			})

			Convey("The filter returned was nil", func() {
				So(filter, ShouldBeNil)
			})
		})
	})

	Convey("Given a nil decoder", t, func() {
		fields := []string{"nop"}
		var decoder GeoIPDecoder = nil

		Convey("When a new geo_ip filter is created", func() {
			filter, err := NewFilterGeoIP(fields, decoder)

			Convey("An error is returned indicating the decoder was nil", func() {
				So(err, ShouldBeError)
			})

			Convey("The filter returned was nil", func() {
				So(filter, ShouldBeNil)
			})
		})
	})
}
