package yaml

// see https://godoc.org/gopkg.in/yaml.v2
import (
	"gopkg.in/yaml.v2"
	"github.com/gookit/config"
)

// Decoder the yaml content decoder
var Decoder config.Decoder = func(blob []byte, ptr interface{}) (err error) {
	return yaml.Unmarshal(blob, ptr)
}

// Encoder the yaml content encoder
var Encoder config.Encoder = func(ptr interface{}) (out []byte, err error) {
	return yaml.Marshal(ptr)
}

// Driver
var Driver = &yamlDriver{config.Yaml}

// yamlDriver for yaml format content
type yamlDriver struct {
	name string
}

// Name
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
