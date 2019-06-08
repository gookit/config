/*
Package hcl is driver use HCL format content as config source

about HCL, please see https://github.com/hashicorp/hcl
*/
package hcl

import (
	"errors"

	"github.com/gookit/config/v2"
	"github.com/hashicorp/hcl"
)

// Decoder the hcl content decoder
var Decoder config.Decoder = hcl.Unmarshal

// Encoder the hcl content encoder
var Encoder config.Encoder = func(ptr interface{}) (out []byte, err error) {
	err = errors.New("HCL: is not support encode data to HCL")
	return
}

// Driver instance for hcl
var Driver = &hclDriver{config.Hcl}

// hclDriver for hcl format content
type hclDriver struct {
	name string
}

// Name
func (d *hclDriver) Name() string {
	return d.name
}

// GetDecoder for hcl
func (d *hclDriver) GetDecoder() config.Decoder {
	return Decoder
}

// GetEncoder for hcl
func (d *hclDriver) GetEncoder() config.Encoder {
	return Encoder
}
