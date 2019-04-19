package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestEnvGetEnvironmentVariable(t *testing.T) {
	Convey("Given an environment variable SOME_VARIABLE is not set", t, func() {
		os.Unsetenv("SOME_VARIABLE")

		Convey("When the variable is retrieved with a fallback of the empty string", func() {
			env := GetEnv("SOME_VARIABLE", "")

			Convey("The empty string should be returned", func() {
				So(env, ShouldEqual, "")
			})
		})

		Convey("When the variable is retrieved with a fallback of \"foo\"", func() {
			env := GetEnv("SOME_VARIABLE", "foo")

			Convey("The string \"foo\" should be returned", func() {
				So(env, ShouldEqual, "foo")
			})
		})
	})

	Convey("Given an environment variable SOME_VARIABLE is set to the empty string", t, func() {
		os.Setenv("SOME_VARIABLE", "")

		Convey("When the variable is retrieved with a fallback of the empty string", func() {
			env := GetEnv("SOME_VARIABLE", "")

			Convey("The empty string should be returned", func() {
				So(env, ShouldEqual, "")
			})
		})
	})

	Convey("Given an environment variable SOME_VARIABLE is set to a non-empty string", t, func() {
		os.Setenv("SOME_VARIABLE", "foo")

		Convey("When the variable is retrieved with a fallback", func() {
			env := GetEnv("SOME_VARIABLE", "bar")

			Convey("The fallback string should not be returned", func() {
				So(env, ShouldNotEqual, "bar")
			})
		})
	})
}
