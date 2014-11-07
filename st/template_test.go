package st

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestShortFormat(t *testing.T) {
	Convey("Given a short format template", t, func() {
		template := `before_$hello__world__argle$_after`

		Convey("When I find string matches with the shortFormat", func() {
			matches := shortFormat.FindAllStringSubmatch(template, -1)

			Convey("Then I expect the string interpolation to be found", func() {
				So(len(matches), ShouldEqual, 1)
				So(matches[0][0], ShouldEqual, `$hello__world__argle$`)
				So(matches[0][1], ShouldEqual, `hello__world__argle`)
			})
		})

		Convey("When I transform the template to a go template", func() {
			text, err := transformShort([]byte(template))

			Convey("Then I expect no errors", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the content should be in go template format", func() {
				So(string(text), ShouldEqual, "before_{{ .hello | world | argle }}_after")
			})
		})
	})
}

func TestLongFormat(t *testing.T) {
	Convey("Given a long format template", t, func() {
		template := `before_$name;format="normalize,lower"$_after`

		Convey("When I find string matches with the longFormat", func() {
			matches := longFormat.FindAllStringSubmatch(template, -1)

			Convey("Then I expect the string interpolation to be found", func() {
				So(len(matches), ShouldEqual, 1)
				So(matches[0][0], ShouldEqual, `$name;format="normalize,lower"$`)
				So(matches[0][1], ShouldEqual, `name;format="normalize,lower"`)
				So(matches[0][2], ShouldEqual, `name`)
				So(matches[0][3], ShouldEqual, `normalize,lower`)
			})
		})

		Convey("When I transform the template to a go template", func() {
			text, err := transformLong([]byte(template))

			Convey("Then I expect no errors", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the content should be in go template format", func() {
				So(string(text), ShouldEqual, "before_{{ .name | normalize | lower }}_after")
			})
		})
	})
}
