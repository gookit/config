package toml

// see https://godoc.org/github.com/BurntSushi/toml
import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/gookit/config"
)

// Decoder the toml content decoder
var Decoder config.Decoder = func(blob []byte, ptr interface{}) (err error) {
	_, err = toml.Decode(string(blob), ptr)

	return
}

// Encoder the toml content encoder
var Encoder config.Encoder = func(ptr interface{}) (out []byte, err error) {
	buf := new(bytes.Buffer)

	err = toml.NewEncoder(buf).Encode(ptr)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}
