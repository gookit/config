/*
Package properties is a driver use Java properties format content as config source

Usage please see readme.

*/
package properties

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/properties"
)

// DriverName string
const DriverName = "properties"

var (
	// Decoder the properties content decoder
	Decoder config.Decoder = properties.Decode

	// Encoder the properties content encoder
	Encoder config.Encoder = properties.Encode

	// Driver for yaml
	Driver = config.NewDriver(DriverName, Decoder, Encoder)
)
