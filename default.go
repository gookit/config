package config

import "encoding/json"

// jsonDecoder
var JsonDecoder Decoder = func(blob []byte, v interface{}) (err error) {
	return json.Unmarshal(blob, v)
}

// a default manager
var defConf = &Config{
	name: DefaultNode,
	format: Yaml,
	decoders: map[string]Decoder{
		Json: JsonDecoder,
	},
	data: make(map[string]interface{}),
}

func Get() {

}

func Set() {

}

func Data() map[string]interface{} {
	return defConf.data
}

func SetOptions(opts *Options)  {
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
