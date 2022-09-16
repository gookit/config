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
var Encoder config.Encoder = parser.Encode

// Driver for ini
var Driver = config.NewDriver(config.Ini, Decoder, Encoder)
