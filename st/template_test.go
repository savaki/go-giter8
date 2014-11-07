package st

import (
	"bytes"
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
			text := transformShort([]byte(template))

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
			text := transformLong([]byte(template))

			Convey("And the content should be in go template format", func() {
				So(string(text), ShouldEqual, "before_{{ .name | normalize | lower }}_after")
			})
		})
	})
}

func TestParse(t *testing.T) {
	Convey("Given the text of a template", t, func() {
		text := []byte(`hello $name;format="lower"$`)

		Convey("When I Parse the template", func() {
			template, err := Parse(text)

			Convey("Then I expect no errors", func() {
				So(err, ShouldBeNil)
			})

			Convey("And I expect a valid template back", func() {
				So(template, ShouldNotBeNil)
			})

			Convey("And I expect the executed template to return the correct value", func() {
				buffer := bytes.NewBuffer([]byte{})
				template.Execute(buffer, map[string]string{"name": "world"})
				So(buffer.String(), ShouldEqual, "hello world")
			})
		})
	})

	Convey("Given an invalid template", t, func() {
		text := []byte("{{ ")

		Convey("When I #Parse the template", func() {
			_, err := Parse(text)

			Convey("Then I expect an error to be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRender(t *testing.T) {
	Convey("Given a text template that uses both the long and short format", t, func() {
		text := []byte(`hello $name;format="lower"$; HELLO $name__upper$`)

		Convey(`When I call #Render`, func() {
			value, err := Render(text, map[string]string{"name": "WoRlD"})

			Convey("Then I expect no errors", func() {
				So(err, ShouldBeNil)
			})

			Convey("And I expect our interpolated string back", func() {
				So(string(value), ShouldEqual, "hello world; HELLO WORLD")
			})
		})
	})
}
