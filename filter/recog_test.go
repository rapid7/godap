package filter

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRecogFilterProcessingBehavior(t *testing.T) {
	Convey("Given a recog filter with a nil mapped field map", t, func() {
		filterRecog, err := NewFilterRecog(nil, "")

		Convey("The recog filter should be nil", func() {
			So(filterRecog, ShouldBeNil)
		})

		Convey("The error should be set (non-nil)", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("Given a recog filter with a mapped field of \"input\" to dns.versionbind", t, func() {
		filterRecog, _ := NewFilterRecog(map[string]string{
			"input": "dns.versionbind",
		}, "")

		Convey("When we process an input doc with a RHEL DNS banner", func() {
			result, _ := filterRecog.Process(map[string]interface{}{
				"input": "9.8.2rc1-RedHat-9.8.2-0.62.rc1.el6_9.2",
			})

			Convey("There should only be one document", func() {
				So(result, ShouldHaveLength, 1)
			})

			Convey("The document should contain recog fields", func() {
				So(result[0]["input.recog.service.product"], ShouldEqual, "BIND")
			})
		})

		Convey("When we process an input doc with no keys matching the filter's mapped fields", func() {
			match_fields := map[string]interface{}{
				"unmatched_field": "9.8.2rc1-RedHat-9.8.2-0.62.rc1.el6_9.2",
			}
			result, _ := filterRecog.Process(match_fields)

			Convey("There should be only one result document", func() {
				So(result, ShouldHaveLength, 1)
			})

			Convey("The result document should equal the input document", func() {
				So(result[0], ShouldResemble, match_fields)
			})
		})
	})
}

func TestRecogFilterDatabaseLoading(t *testing.T) {
	Convey("Given a new recog filter created with an invalid database path", t, func() {
		filterRecog, err := NewFilterRecog(map[string]string{"some": "field"}, "/some/invalid/path")

		Convey("The recog filter should be nil", func() {
			So(filterRecog, ShouldEqual, nil)
		})

		Convey("The error returned by NewFilterRecog should be set", func() {
			So(err, ShouldBeError)
		})
	})

	Convey("Given a new recog filter created with an empty database path", t, func() {
		filterRecog, err := NewFilterRecog(map[string]string{"some": "field"}, "")

		Convey("The recog filter should be non-nil", func() {
			So(filterRecog, ShouldNotBeNil)
		})

		Convey("The error returned by NewFilterRecog should be unset (nil)", func() {
			So(err, ShouldBeNil)
		})
	})
}
