package template

import (
	"regexp"
)

var (
	stringEscapeRe = regexp.MustCompile(`\\\$`)
)

func giter8ify(text []byte) []byte {
	return stringEscapeRe.ReplaceAll(text, []byte(`$`))
}
