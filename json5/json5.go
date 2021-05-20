// Package json5 use the https://github.com/yosuke-furukawa/json5 for parse json5
package json5

import (
	"github.com/gookit/config/v2"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

const NAME = "json5"

var (
	// Decoder for json
	Decoder config.Decoder = json5.Unmarshal

	// Encoder for json5
	Encoder config.Encoder = json5.Marshal

	// Driver for json5
	Driver = &json5Driver{name: NAME}
)

// json5Driver for json5 format content
type json5Driver struct {
	name string
}

// Name get name
func (d *json5Driver) Name() string {
	return d.name
}

// GetDecoder for json5
func (d *json5Driver) GetDecoder() config.Decoder {
	return Decoder
}

// GetEncoder for json5
func (d *json5Driver) GetEncoder() config.Encoder {
	return Encoder
}
