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

// Driver for toml
var Driver = &tomlDriver{config.Toml}

// tomlDriver for toml format content
type tomlDriver struct {
	name string
}

// Name get name
func (d *tomlDriver) Name() string {
	return d.name
}

// GetDecoder for toml
func (d *tomlDriver) GetDecoder() config.Decoder {
	return Decoder
}

// GetEncoder for toml
func (d *tomlDriver) GetEncoder() config.Encoder {
	return Encoder
}
