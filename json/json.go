// Package json use the https://github.com/json-iterator/go for parse json
package json

import (
	"github.com/goccy/go-json"
	"github.com/gookit/config/v2"
	"github.com/gookit/goutil/jsonutil"
)

var (
	// Decoder for json
	Decoder config.Decoder = func(data []byte, v any) (err error) {
		if config.JSONAllowComments {
			str := jsonutil.StripComments(string(data))
			return json.Unmarshal([]byte(str), v)
		}
		return json.Unmarshal(data, v)
	}

	// Encoder for json
	Encoder config.Encoder = json.Marshal
	// Driver for json
	Driver = config.NewDriver(config.JSON, Decoder, Encoder)
)
