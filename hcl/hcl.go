/*
Package hcl is driver use HCL format content as config source

about HCL, please see https://github.com/hashicorp/hcl

NOTE: Deprecated, The HCL driver will no longer be built-in, please add the HCL driver yourself

*/
package hcl

import (
	"errors"

	"github.com/gookit/config/v2"
	// "github.com/hashicorp/hcl"
)

// Decoder the hcl content decoder
// var Decoder config.Decoder = hcl.Unmarshal
var Decoder config.Decoder = func(blob []byte, v any) (err error) {
	return errors.New("HCL: is not support decode data from HCL")
}

// Encoder the hcl content encoder
var Encoder config.Encoder = func(ptr any) (out []byte, err error) {
	err = errors.New("HCL: is not support encode data to HCL")
	return
}

// Driver instance for hcl
var Driver = config.NewDriver(config.Hcl, Decoder, Encoder).WithAlias("conf")
