package json

// use the https://github.com/json-iterator/go
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
