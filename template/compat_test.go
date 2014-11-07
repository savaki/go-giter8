package template

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGiter8ify(t *testing.T) {
	Convey("Given a string with escaped $", t, func() {
		text := []byte(`-> \$ <-`)

		Convey("When I invoke #giter8ify", func() {
			result := giter8ify(text)

			Convey("Then I expect the escaped $ to be unescaped", func() {
				So(string(result), ShouldEqual, `-> $ <-`)
			})
		})
	})
}
