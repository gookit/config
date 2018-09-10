package config

// default json driver(encoder/decoder)
import (
	"encoding/json"
	"regexp"
	"strings"
)

// JSONDecoder for json decode
var JSONDecoder Decoder = func(data []byte, v interface{}) (err error) {
	str := ClearJSONComments(string(data))
	return json.Unmarshal([]byte(str), v)
}

// JSONEncoder for json encode
var JSONEncoder Encoder = json.Marshal

// JSONDriver instance fot json
var JSONDriver = &jsonDriver{JSON}

// jsonDriver for json format content
type jsonDriver struct {
	name string
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
var jsonSLComments = regexp.MustCompile(`(?s://.*?)[\r\n]`)
var jsonMLComments = regexp.MustCompile(`(?s:/\*.*?\*/\s*)`)

// ClearJSONComments for a JSON string
func ClearJSONComments(str string) string {
	if !strings.Contains(str, "//") && !strings.Contains(str, "/*") {
		return str
	}

	str = jsonMLComments.ReplaceAllString(str, "")
	str = jsonSLComments.ReplaceAllString(str, "")

	return strings.TrimSpace(str)
}
