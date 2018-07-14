package config

import "encoding/json"

// jsonDecoder
var JsonDecoder Decoder = func(blob []byte, v interface{}) (err error) {
	return json.Unmarshal(blob, v)
}

// JsonEncoder
var JsonEncoder Encoder = func(v interface{}) (out []byte, err error) {
	return json.Marshal(v)
}
