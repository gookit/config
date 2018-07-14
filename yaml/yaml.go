package yaml

import (
	"gopkg.in/yaml.v2"
	"github.com/gookit/config"
)

// Decoder
var Decoder config.Decoder = func(blob []byte, ptr interface{}) (err error) {
	return yaml.Unmarshal(blob, ptr)
}

// Encode
func Encode(ptr interface{}) (out []byte, err error) {
	return yaml.Marshal(ptr)
}
