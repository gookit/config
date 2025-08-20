/*
Package hclv2 is driver use HCL format content as config source

about HCL, please see https://github.com/hashicorp/hcl
docs for HCL v2 https://pkg.go.dev/github.com/hashicorp/hcl/v2

NOTE: Deprecated, The HCLv2 driver will no longer be built-in, please add the HCLv2 driver yourself
*/
package hclv2

import (
	"errors"

	"github.com/nstrlabs/config"
	// "github.com/hashicorp/hcl/v2/hclsimple"
	// "github.com/hashicorp/hcl/v2/hclparse"
)

// Decoder the hcl content decoder
var Decoder config.Decoder = func(blob []byte, v any) (err error) {
	// hclparse.NewParser().ParseHCL(blob, "hcl2/config.hcl")
	// TODO hcl2 decode data to map ptr will report error
	// return hclsimple.Decode("hcl2/config.hcl", blob, nil, v)
	return errors.New("HCLv2: is not support decode data from HCL")
}

// Encoder the hcl content encoder
var Encoder config.Encoder = func(ptr any) (out []byte, err error) {
	err = errors.New("HCLv2: is not support encode data to HCL")
	return
}

// Driver instance for hcl
var Driver = config.NewDriver(config.Hcl, Decoder, Encoder).WithAlias("conf")
