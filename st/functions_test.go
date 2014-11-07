package st

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCap(t *testing.T) {
	Convey("When I #Capitalize a string", t, func() {
		result := Capitalize("hello")

		Convey("Then only the first letter is capitalized", func() {
			So(result, ShouldEqual, "Hello")
		})
	})

	Convey("I can #Capitalize a zero lengthed string", t, func() {
		So(Capitalize(""), ShouldEqual, "")
	})

	Convey("I can #Capitalize a string of length 1", t, func() {
		So(Capitalize("h"), ShouldEqual, "H")
	})
}

func TestRandom(t *testing.T) {
	Convey("When I #Random a string", t, func() {
		result := Random("hello")
		So(len(result), ShouldBeGreaterThan, len("hello"))
	})
}

func TestStart(t *testing.T) {
	Convey("When I #Start a string", t, func() {
		result := Start("hello world")
		So(result, ShouldEqual, "Hello World")
	})
}

func TestNormalize(t *testing.T) {
	Convey("When I #Normalize a string", t, func() {
		result := Normalize("The Time Has Come")
		So(result, ShouldEqual, "the-time-has-come")
	})
}

func TestCamel(t *testing.T) {
	Convey("When I #Camel a string", t, func() {
		result := Camel("hello world")
		So(result, ShouldEqual, "HelloWorld")
	})
}

func TestCamelLower(t *testing.T) {
	Convey("When I #CamelLower a string", t, func() {
		result := CamelLower("hello world")
		So(result, ShouldEqual, "helloWorld")
	})
}

func TestSnake(t *testing.T) {
	Convey("When I #Snake a string", t, func() {
		result := Snake("hello world.argle.bargle")
		So(result, ShouldEqual, "hello_world_argle_bargle")
	})
}
