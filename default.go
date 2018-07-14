package config

// this a default config manager
var defConf = New("default")

// Get
func Get(key string, findByPath ...bool) (value interface{}, ok bool) {
	return defConf.Get(key, findByPath...)
}

// GetStringMap get config data as a map[string]string
func GetStringMap(key string) (mp map[string]string, ok bool) {
	return defConf.GetStringMap(key)
}

// GetStringArr
func GetStringArr(key string) (arr []string, ok bool) {
	return defConf.GetStringArr(key)
}

func Set() {

}

// Data
func Data() map[string]interface{} {
	return defConf.data
}

func SetOptions(opts *Options) {
	defConf.options = opts
}

// SetDecoder
func SetDecoder(format string, decoder Decoder) {
	defConf.SetDecoder(format, decoder)
}

// LoadFiles
func LoadFiles(sourceFiles ...string) (err error) {
	return defConf.LoadFiles(sourceFiles...)
}

// LoadData
func LoadData(dataSource ...interface{}) (err error) {
	return defConf.LoadData(dataSource...)
}

// LoadSources
func LoadSources(format string, sourceCode ...[]byte) (err error) {
	return defConf.LoadSources(format, sourceCode...)
}
