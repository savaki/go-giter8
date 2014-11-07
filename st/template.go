package st

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"
)

var (
	// $name__filter1__filter2$
	shortFormat = regexp.MustCompile(`\$([a-zA-Z0-9]+(__[a-zA-Z0-9]+)*)\$`)

	// $name;format="filter1,filter2"$
	longFormat = regexp.MustCompile(`\$((\w+);format="([^"]+)")\$`)
)

// converts data on the other end of the reader to a golang template
func Parse(text []byte) (*template.Template, error) {
	text, err := transform(text)
	if err != nil {
		return nil, err
	}

	return template.New("template").Funcs(funcMap).Parse(string(text))
}

func Render(text []byte, data interface{}) ([]byte, error) {
	t, err := Parse(text)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer([]byte{})
	err = t.Execute(buffer, data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// transform templates in the short format to go templates e.g. name__filter1__filter2
func transformShort(text []byte) ([]byte, error) {
	matches := shortFormat.FindAllSubmatch(text, -1)
	for _, match := range matches {
		segments := bytes.Split(match[1], []byte("__"))
		macro := fmt.Sprintf("{{ .%s }}", bytes.Join(segments, []byte(" | ")))
		text = bytes.Replace(text, match[0], []byte(macro), -1)
	}

	return text, nil
}

// transform templates in the long format to go templates e.g. name;filter="lower,snake"
func transformLong(text []byte) ([]byte, error) {
	matches := longFormat.FindAllSubmatch(text, -1)
	for _, match := range matches {
		field := match[2]
		filters := match[3]

		segments := [][]byte{field}
		segments = append(segments, bytes.Split(filters, []byte(","))...)

		macro := fmt.Sprintf("{{ .%s }}", bytes.Join(segments, []byte(" | ")))
		text = bytes.Replace(text, match[0], []byte(macro), -1)
	}

	return text, nil
}

// helper to combine both short and long transforms
func transform(text []byte) ([]byte, error) {
	text, err := transformShort(text)
	if err != nil {
		return nil, err
	}

	text, err = transformLong(text)
	if err != nil {
		return nil, err
	}

	return text, nil
}
