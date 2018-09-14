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
func Get(key string, findByPath ...bool) (value interface{}, ok bool) {
	return dc.Get(key, findByPath...)
}

// Int get a int by key
func Int(key string) (value int, ok bool) {
	return dc.Int(key)
}

// DefInt get a int value, if not found return default value
func DefInt(key string, def int) int {
	return dc.DefInt(key, def)
}

// Int64 get a int64 by key
func Int64(key string) (value int64, ok bool) {
	return dc.Int64(key)
}

// DefInt64 get a int64 with a default value
func DefInt64(key string, def int64) int64 {
	return dc.DefInt64(key, def)
}

// Bool get a bool by key
func Bool(key string) (value bool, ok bool) {
	return dc.Bool(key)
}

// DefBool get a bool value, if not found return default value
func DefBool(key string, def bool) bool {
	return dc.DefBool(key, def)
}

// Float get a bool by key
func Float(key string) (value float64, ok bool) {
	return dc.Float(key)
}

// DefFloat get a float64 value, if not found return default value
func DefFloat(key string, def float64) float64 {
	return dc.DefFloat(key, def)
}

// String get a string by key
func String(key string) (value string, ok bool) {
	return dc.String(key)
}

// DefString get a string value, if not found return default value
func DefString(key string, def string) string {
	return dc.DefString(key, def)
}

// Ints  get config data as a int slice/array
func Ints(key string) (arr []int, ok bool) {
	return dc.Ints(key)
}

// IntMap get config data as a map[string]int
func IntMap(key string) (mp map[string]int, ok bool) {
	return dc.IntMap(key)
}

// Strings get strings by key
func Strings(key string) (arr []string, ok bool) {
	return dc.Strings(key)
}

// StringMap get config data as a map[string]string
func StringMap(key string) (mp map[string]string, ok bool) {
	return dc.StringMap(key)
}

/*************************************************************
 * load config data to default instance
 *************************************************************/

// LoadFiles load one or multi files
func LoadFiles(sourceFiles ...string) (err error) {
	return dc.LoadFiles(sourceFiles...)
}

// LoadExists load one or multi files, will ignore not exist
func LoadExists(sourceFiles ...string) (err error) {
	return dc.LoadExists(sourceFiles...)
}

// LoadData load one or multi data
func LoadData(dataSource ...interface{}) (err error) {
	return dc.LoadData(dataSource...)
}

// LoadSources load one or multi byte data
func LoadSources(format string, src []byte, more ...[]byte) (err error) {
	return dc.LoadSources(format, src, more...)
}

// LoadStrings load one or multi string
func LoadStrings(format string, str string, more ...string) (err error) {
	return dc.LoadStrings(format, str, more...)
}

/*************************************************************
 * helper functions for the default instance
 *************************************************************/

// Set val by key
func Set(key string, val interface{}, setByPath ...bool) (err error) {
	return dc.Set(key, val, setByPath...)
}

// WriteTo a writer
func WriteTo(out io.Writer) (n int64, err error) {
	return dc.WriteTo(out)
}

// DumpTo a writer and use format
func DumpTo(out io.Writer, format string) (n int64, err error) {
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
