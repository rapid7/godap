package filter

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRecogFilterProcessingBehavior(t *testing.T) {
	Convey("Given a recog filter with a mapped field of \"input\" to dns.versionbind", t, func() {
		filterRecog, _ := NewFilterRecog(map[string]string{
			"input": "dns.versionbind",
		}, "")

		Convey("When we process an input doc with a RHEL DNS banner", func() {
			result, _ := filterRecog.Process(map[string]interface{}{
				"input": "9.8.2rc1-RedHat-9.8.2-0.62.rc1.el6_9.2",
			})

			Convey("The document should contain recog fields", func() {
				So(result[0]["input.recog.service.product"], ShouldEqual, "BIND")
			})
		})

		Convey("When we process an input doc with no keys matching the filter's mapped fields", func() {
			result, _ := filterRecog.Process(map[string]interface{}{
				"unmatched_field": "9.8.2rc1-RedHat-9.8.2-0.62.rc1.el6_9.2",
			})

			Convey("There should be no recog fields in the result document", func() {
				So(result[0]["input.recog.service.product"], ShouldEqual, nil)
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

		Convey("The error returned by NewFilterRecog should be nil", func() {
			So(err, ShouldBeNil)
		})
	})
}
