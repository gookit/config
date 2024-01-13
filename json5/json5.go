// Package json5 support for parse and load json5
package json5

import (
	"encoding/json"

	"github.com/gookit/config/v2"
	"github.com/titanous/json5"
)

// Name for driver
const Name = "json5"

// NAME for driver
const NAME = Name

// JSONMarshalIndent if not empty, will use json.MarshalIndent for encode data.
var JSONMarshalIndent string

var (
	// Decoder for json
	Decoder config.Decoder = json5.Unmarshal

	// Encoder for json5
	Encoder config.Encoder = func(v any) (out []byte, err error) {
		if len(JSONMarshalIndent) == 0 {
			return json.Marshal(v)
		}
		return json.MarshalIndent(v, "", JSONMarshalIndent)
	}

	// Driver for json5
	Driver = config.NewDriver(Name, Decoder, Encoder)
)
