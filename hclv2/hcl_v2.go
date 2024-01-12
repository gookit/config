/*
Package hclv2 is driver use HCL format content as config source

about HCL, please see https://github.com/hashicorp/hcl
docs for HCL v2 https://pkg.go.dev/github.com/hashicorp/hcl/v2
*/
package hclv2

import (
	"errors"

	"github.com/gookit/config/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

// Decoder the hcl content decoder
var Decoder config.Decoder = func(blob []byte, v any) (err error) {
	// TODO hcl2 decode data to map ptr will report error
	return hclsimple.Decode("hcl2/config.hcl", blob, nil, v)
}

// Encoder the hcl content encoder
var Encoder config.Encoder = func(ptr any) (out []byte, err error) {
	err = errors.New("HCLv2: is not support encode data to HCL")
	return
}

// Driver instance for hcl
var Driver = config.NewDriver(config.Hcl, Decoder, Encoder).WithAlias("conf")
