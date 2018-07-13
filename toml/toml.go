package toml

// https://godoc.org/github.com/BurntSushi/toml
import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/gookit/config"
)

// Decoder
var Decoder config.Decoder = func (blob []byte, ptr interface{}) (err error) {
	_, err = toml.Decode(string(blob), ptr)

	return
}

// Encode
func Encode(ptr interface{}) (str string, err error) {
	buf := new(bytes.Buffer)
	err = toml.NewEncoder(buf).Encode(ptr)

	if err != nil {
	    return
	}

	str = buf.String()
	return
}
