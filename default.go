package config

// this is a default config manager
var dc = New("default")

// Get
func Get(key string, findByPath ...bool) (value interface{}, ok bool) {
	return dc.Get(key, findByPath...)
}

// GetInt
func GetInt(key string) (value int, ok bool) {
	return dc.GetInt(key)
}

// DefInt get a int value, if not found return default value
func DefInt(key string, def int) int {
	return dc.DefInt(key, def)
}

// GetBool
func GetBool(key string) (value bool, ok bool) {
	return dc.GetBool(key)
}

// DefBool get a bool value, if not found return default value
func DefBool(key string, def bool) bool {
	return dc.DefBool(key, def)
}

// GetString
func GetString(key string) (value string, ok bool) {
	return dc.GetString(key)
}

// DefString get a string value, if not found return default value
func DefString(key string, def string) string {
	return dc.DefString(key, def)
}

// GetStringMap get config data as a map[string]string
func GetStringMap(key string) (mp map[string]string, ok bool) {
	return dc.GetStringMap(key)
}

// GetStringArr
func GetStringArr(key string) (arr []string, ok bool) {
	return dc.GetStringArr(key)
}

func Set() {

}

// Data
func Data() map[string]interface{} {
	return dc.data
}

func SetOptions(opts *Options) {
	dc.options = opts
}

// SetDecoder
func SetDecoder(format string, decoder Decoder) {
	dc.SetDecoder(format, decoder)
}

// LoadFiles
func LoadFiles(sourceFiles ...string) (err error) {
	return dc.LoadFiles(sourceFiles...)
}

// LoadData
func LoadData(dataSource ...interface{}) (err error) {
	return dc.LoadData(dataSource...)
}

// LoadSources
func LoadSources(format string, sourceCode ...[]byte) (err error) {
	return dc.LoadSources(format, sourceCode...)
}
