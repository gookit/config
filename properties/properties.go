/*
Package properties is a driver use Java properties format content as config source

Usage please see readme.

*/
package properties

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/properties"
)

// Name string
const Name = "properties"

var (
	// Decoder the properties content decoder
	Decoder config.Decoder = properties.Decode

	// Encoder the properties content encoder
	Encoder config.Encoder = properties.Encode

	// Driver for properties
	Driver = config.NewDriver(Name, Decoder, Encoder)
)
