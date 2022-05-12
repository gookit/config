// Package json use the https://github.com/json-iterator/go for parse json
package json

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/goutil/jsonutil"
	"github.com/json-iterator/go"
)

var parser = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	// Decoder for json
	Decoder config.Decoder = func(data []byte, v interface{}) (err error) {
		if config.JSONAllowComments {
			str := jsonutil.StripComments(string(data))
			return parser.Unmarshal([]byte(str), v)
		}
		return parser.Unmarshal(data, v)
	}

	// Encoder for json
	Encoder config.Encoder = parser.Marshal
	// Driver for json
	Driver = config.NewDriver(config.JSON, Decoder, Encoder)
)
