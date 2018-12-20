package config

// default json driver(encoder/decoder)
import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	"text/scanner"
)

// JSONDecoder for json decode
var JSONDecoder Decoder = func(data []byte, v interface{}) (err error) {
	str := StripJSONComments(string(data))
	return json.Unmarshal([]byte(str), v)
}

// JSONEncoder for json encode
var JSONEncoder Encoder = json.Marshal

// JSONDriver instance fot json
var JSONDriver = &jsonDriver{name: JSON}

// jsonDriver for json format content
type jsonDriver struct {
	name          string
	ClearComments bool
}

// Name
func (d *jsonDriver) Name() string {
	return d.name
}

// GetDecoder for json
func (d *jsonDriver) GetDecoder() Decoder {
	return JSONDecoder
}

// GetEncoder for json
func (d *jsonDriver) GetEncoder() Encoder {
	return JSONEncoder
}

// `(?s:` enable match multi line
var jsonMLComments = regexp.MustCompile(`(?s:/\*.*?\*/\s*)`)

// StripJSONComments for a JSON string
func StripJSONComments(src string) string {
	// multi line comments
	if strings.Contains(src, "/*") {
		src = jsonMLComments.ReplaceAllString(src, "")
	}

	// single line comments
	if !strings.Contains(src, "//") {
		return strings.TrimSpace(src)
	}

	var s scanner.Scanner

	s.Init(strings.NewReader(src))
	s.Filename = "comments"
	s.Mode ^= scanner.SkipComments // don't skip comments

	buf := new(bytes.Buffer)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		txt := s.TokenText()
		if !strings.HasPrefix(txt, "//") && !strings.HasPrefix(txt, "/*") {
			buf.WriteString(txt)
			// } else {
			// fmt.Printf("%s: %s\n", s.Position, txt)
		}
	}
	return buf.String()
}
