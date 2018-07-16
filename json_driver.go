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

// JsonDriver for json format content
type JsonDriver struct{}

// GetDecoder for json
func (d *JsonDriver) GetDecoder() Decoder {
	return JsonDecoder
}

// GetEncoder for json
func (d *JsonDriver) GetEncoder() Encoder {
	return JsonEncoder
}
