/*
Package other is an example of a custom driver
*/
package other

import (
	"github.com/nstrlabs/config"
	"github.com/nstrlabs/config/ini"
)

// DriverName string
const DriverName = "other"

var (
	// Encoder is the encoder for this driver
	Encoder = ini.Encoder
	// Decoder is the decoder for this driver
	Decoder = ini.Decoder
	// Driver is the exported symbol
	Driver = config.NewDriver(DriverName, Decoder, Encoder)
)
