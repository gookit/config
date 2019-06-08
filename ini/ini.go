/*
Package ini is driver use INI format content as config source

about ini parse, please see https://github.com/gookit/ini/parser
*/
package ini

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/ini/v2/parser"
)

// Decoder the ini content decoder
var Decoder config.Decoder = parser.Decode

// Encoder encode data to ini content
var Encoder config.Encoder = func(ptr interface{}) (out []byte, err error) {
	return parser.Encode(ptr)
}

// Driver for ini
var Driver = &iniDriver{config.Ini}

// iniDriver for ini format content
type iniDriver struct {
	name string
}

// Name get name
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
