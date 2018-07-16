// use the https://github.com/json-iterator/go for parse json
package json

import (
	"github.com/json-iterator/go"
	"github.com/gookit/config"
)

// Decoder
var Decoder config.Decoder = func(blob []byte, v interface{}) (err error) {
	var parser = jsoniter.ConfigCompatibleWithStandardLibrary
	return parser.Unmarshal(blob, v)
}

// Encoder
var Encoder config.Encoder = func(v interface{}) (out []byte, err error) {
	var parser = jsoniter.ConfigCompatibleWithStandardLibrary
	return parser.Marshal(v)
}

// Driver
var Driver = &jsonDriver{config.Json}

// jsonDriver for json format content
type jsonDriver struct {
	name string
}

// Name
func (d *jsonDriver) Name() string {
	return d.name
}

// GetDecoder for json
func (d *jsonDriver) GetDecoder() config.Decoder {
	return Decoder
}

// GetEncoder for json
func (d *jsonDriver) GetEncoder() config.Encoder {
	return Encoder
}
