package config

import "io"

// This is a default config manager instance
var dc = New("default")

// Default get the default instance
func Default() *Config {
	return dc
}

/*************************************************************
 * config setting for the default instance
 *************************************************************/

// GetOptions get options
func GetOptions() *Options {
	return dc.Options()
}

// WithOptions with options
func WithOptions(opts ...func(*Options)) {
	dc.WithOptions(opts...)
}

// AddDriver set a decoder and encoder driver for a format.
func AddDriver(driver Driver) {
	dc.AddDriver(driver)
}

// SetDecoder add/set a format decoder
func SetDecoder(format string, decoder Decoder) {
	dc.SetDecoder(format, decoder)
}

// SetEncoder set a encoder for the format
func SetEncoder(format string, encoder Encoder) {
	dc.SetEncoder(format, encoder)
}

/*************************************************************
 * read config data
 *************************************************************/

// Get get a value by key
func Get(key string, findByPath ...bool) (interface{}, bool) {
	return dc.Get(key, findByPath...)
}

// Int get a int by key
func Int(key string) (int, bool) {
	return dc.Int(key)
}

// DefInt get a int value, if not found return default value
func DefInt(key string, defVal ...int) int {
	return dc.DefInt(key, defVal...)
}

// Int64 get a int64 by key
func Int64(key string) (int64, bool) {
	return dc.Int64(key)
}

// DefInt64 get a int64 with a default value
func DefInt64(key string, defVal ...int64) int64 {
	return dc.DefInt64(key, defVal...)
}

// Bool get a bool by key
func Bool(key string) (bool, bool) {
	return dc.Bool(key)
}

// DefBool get a bool value, if not found return default value
func DefBool(key string, defVal ...bool) bool {
	return dc.DefBool(key, defVal...)
}

// Float get a bool by key
func Float(key string) (float64, bool) {
	return dc.Float(key)
}

// DefFloat get a float64 value, if not found return default value
func DefFloat(key string, defVal ...float64) float64 {
	return dc.DefFloat(key, defVal...)
}

// String get a string by key
func String(key string) (string, bool) {
	return dc.String(key)
}

// DefString get a string value, if not found return default value
func DefString(key string, defVal ...string) string {
	return dc.DefString(key, defVal...)
}

// Ints  get config data as a int slice/array
func Ints(key string) ([]int, bool) {
	return dc.Ints(key)
}

// IntMap get config data as a map[string]int
func IntMap(key string) (map[string]int, bool) {
	return dc.IntMap(key)
}

// Strings get strings by key
func Strings(key string) ([]string, bool) {
	return dc.Strings(key)
}

// StringMap get config data as a map[string]string
func StringMap(key string) (map[string]string, bool) {
	return dc.StringMap(key)
}

/*************************************************************
 * load config data to default instance
 *************************************************************/

// LoadFlags load data from cli flags
func LoadFlags(keys []string) error {
	return dc.LoadFlags(keys)
}

// LoadFiles load one or multi files
func LoadFiles(sourceFiles ...string) error {
	return dc.LoadFiles(sourceFiles...)
}

// LoadExists load one or multi files, will ignore not exist
func LoadExists(sourceFiles ...string) error {
	return dc.LoadExists(sourceFiles...)
}

// LoadData load one or multi data
func LoadData(dataSource ...interface{}) error {
	return dc.LoadData(dataSource...)
}

// LoadSources load one or multi byte data
func LoadSources(format string, src []byte, more ...[]byte) error {
	return dc.LoadSources(format, src, more...)
}

// LoadStrings load one or multi string
func LoadStrings(format string, str string, more ...string) error {
	return dc.LoadStrings(format, str, more...)
}

/*************************************************************
 * helper functions for the default instance
 *************************************************************/

// Set val by key
func Set(key string, val interface{}, setByPath ...bool) error {
	return dc.Set(key, val, setByPath...)
}

// WriteTo a writer
func WriteTo(out io.Writer) (n int64, err error) {
	return dc.WriteTo(out)
}

// DumpTo a writer and use format
func DumpTo(out io.Writer, format string) (int64, error) {
	return dc.DumpTo(out, format)
}

// Data return all config data
func Data() map[string]interface{} {
	return dc.Data()
}

// ClearAll data and caches
func ClearAll() {
	dc.ClearAll()
}
