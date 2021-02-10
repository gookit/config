/*
Package yaml is a driver use YAML format content as config source

Usage please see example:

*/
package yamlv3

// see https://pkg.go.dev/gopkg.in/yaml.v3
import (
	"github.com/gookit/config/v2"
	"gopkg.in/yaml.v3"
)

// Decoder the yaml content decoder
var Decoder config.Decoder = yaml.Unmarshal

// Encoder the yaml content encoder
var Encoder config.Encoder = yaml.Marshal

// Driver for yaml
var Driver = &yamlDriver{config.Yaml}

// yamlDriver for yaml format content
type yamlDriver struct {
	name string
}

// Name for driver
func (d *yamlDriver) Name() string {
	return d.name
}

// GetDecoder for yaml
func (d *yamlDriver) GetDecoder() config.Decoder {
	return Decoder
}

// GetEncoder for yaml
func (d *yamlDriver) GetEncoder() config.Encoder {
	return Encoder
}
