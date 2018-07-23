package config

// default json driver(encoder/decoder)
import "encoding/json"

// JsonDecoder
var JsonDecoder Decoder = json.Unmarshal

// JsonEncoder
var JsonEncoder Encoder = json.Marshal

// var JsonEncoder Encoder = func(v interface{}) (out []byte, err error) {
// 	return json.Marshal(v)
// }

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
