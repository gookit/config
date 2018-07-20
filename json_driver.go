package config

// default json driver(encoder/decoder)
import "encoding/json"

// JsonDecoder
var JsonDecoder Decoder = func(blob []byte, v interface{}) (err error) {
	return json.Unmarshal(blob, v)
}

// JsonEncoder
var JsonEncoder Encoder = func(v interface{}) (out []byte, err error) {
	return json.Marshal(v)
}

// JsonDriver
var JsonDriver = &jsonDriver{Json}

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
	return JsonDecoder
}

// GetEncoder for json
func (d *jsonDriver) GetEncoder() Encoder {
	return JsonEncoder
}
