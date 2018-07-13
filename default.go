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

// a default manager
var defConf = &Config{
	name: "default",
	data: make(map[string]interface{}),

	defFormat: Json,

	encoders: map[string]Encoder{Json: JsonEncoder},
	decoders: map[string]Decoder{Json: JsonDecoder},
}

func Get() {

}

func Set() {

}

func Data() map[string]interface{} {
	return defConf.data
}

func SetOptions(opts *Options) {
	defConf.options = opts
}

func SetDecoder(format string, decoder Decoder) {
	defConf.SetDecoder(format, decoder)
}

func LoadFiles(sourceFiles ...string) (err error) {
	return defConf.LoadFiles(sourceFiles...)
}

func LoadData(dataSource ...interface{}) (err error) {
	return defConf.LoadData(dataSource...)
}

func LoadSources(format string, sourceCode ...[]byte) (err error) {
	return defConf.LoadSources(format, sourceCode...)
}
