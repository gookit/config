package config

// default json driver(encoder/decoder)
import "encoding/json"

// JSONDecoder
var JSONDecoder Decoder = json.Unmarshal

// JSONEncoder
var JSONEncoder Encoder = json.Marshal

// var JSONEncoder Encoder = func(v interface{}) (out []byte, err error) {
// 	return json.Marshal(v)
// }

// JSONDriver
var JSONDriver = &jsonDriver{JSON}

// jsonDriver for json format content
type jsonDriver struct {
	name string
}

// Name
func (d *jsonDriver) Name() string {
	return d.name
}

// GetDecoder for json
func (d *jsonDriver) GetDecoder() Decoder {
	return JSONDecoder
}

// GetEncoder for json
func (d *jsonDriver) GetEncoder() Encoder {
	return JSONEncoder
}
