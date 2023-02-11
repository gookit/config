/*
Package toml is driver use TOML format content as config source

Usage please see example:
*/
package toml

// see https://godoc.org/github.com/BurntSushi/toml
import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/gookit/config/v2"
)

// Decoder the toml content decoder
var Decoder config.Decoder = func(blob []byte, ptr any) (err error) {
	_, err = toml.Decode(string(blob), ptr)
	return
}

// Encoder the toml content encoder
var Encoder config.Encoder = func(ptr any) (out []byte, err error) {
	buf := new(bytes.Buffer)

	err = toml.NewEncoder(buf).Encode(ptr)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}

// Driver for toml format
var Driver = config.NewDriver(config.Toml, Decoder, Encoder)
