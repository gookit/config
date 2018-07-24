// use the https://github.com/json-iterator/go for parse json
package json

import (
	"github.com/gookit/config"
	"github.com/json-iterator/go"
)

var parser = jsoniter.ConfigCompatibleWithStandardLibrary

// Decoder for json
var Decoder config.Decoder = parser.Unmarshal

// Encoder for json
var Encoder config.Encoder = parser.Marshal

// Driver for json
var Driver = &jsonDriver{config.JSON}

// jsonDriver for json format content
type jsonDriver struct {
	name string
}

// Name get name
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
