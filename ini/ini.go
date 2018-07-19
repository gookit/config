/*
This is driver use INI format content as config source

about ini parse, please see https://github.com/gookit/ini/parser
 */
package ini

import (
	"github.com/gookit/config"
	"github.com/gookit/ini/parser"
	"errors"
)

// Decoder the ini content decoder
var Decoder config.Decoder = func(blob []byte, ptr interface{}) (err error) {
	return parser.Decode(blob, ptr)
}

// Encoder encode data to ini content
var Encoder config.Encoder = func(ptr interface{}) (out []byte, err error) {
	err = errors.New("INI: is not support encode data to INI")
	return
}

// Driver
var Driver = &iniDriver{config.Ini}

// iniDriver for ini format content
type iniDriver struct {
	name string
}

// Name
func (d *iniDriver) Name() string {
	return d.name
}

// GetDecoder for ini
func (d *iniDriver) GetDecoder() config.Decoder {
	return Decoder
}

// GetEncoder for ini
func (d *iniDriver) GetEncoder() config.Encoder {
	return Encoder
}
