package config

import "io"

// This is a default config manager instance
var dc = New("default")

/*************************************************************
 * config setting for the default instance
 *************************************************************/

// SetOptions
func SetOptions(opts *Options) {
	dc.SetOptions(opts)
}

// AddDriver set a decoder and encoder driver for a format.
func AddDriver(format string, driver Driver) {
	dc.AddDriver(format, driver)
}

// DecoderEncoder set a decoder and encoder for a format.
func DecoderEncoder(format string, decoder Decoder, encoder Encoder) {
	dc.DecoderEncoder(format, decoder, encoder)
}

// SetDecoder add/set a format decoder
func SetDecoder(format string, decoder Decoder) {
	dc.SetDecoder(format, decoder)
}

// SetDecoders
func SetDecoders(decoders map[string]Decoder) {
	dc.SetDecoders(decoders)
}

// SetEncoder
func SetEncoder(format string, encoder Encoder) {
	dc.SetEncoder(format, encoder)
}

// SetEncoders
func SetEncoders(encoders map[string]Encoder) {
	dc.SetEncoders(encoders)
}

/*************************************************************
 * read config data
 *************************************************************/

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

// GetInt64
func GetInt64(key string) (value int64, ok bool) {
	return dc.GetInt64(key)
}

// DefInt64
func DefInt64(key string, def int64) int64 {
	return dc.DefInt64(key, def)
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

// GetIntArr  get config data as a int slice/array
func GetIntArr(key string) (arr []int, ok bool) {
	return dc.GetIntArr(key)
}

// GetIntMap get config data as a map[string]int
func GetIntMap(key string) (mp map[string]int, ok bool) {
	return dc.GetIntMap(key)
}

// GetStringArr
func GetStringArr(key string) (arr []string, ok bool) {
	return dc.GetStringArr(key)
}

// GetStringMap get config data as a map[string]string
func GetStringMap(key string) (mp map[string]string, ok bool) {
	return dc.GetStringMap(key)
}

/*************************************************************
 * load config data to default instance
 *************************************************************/

// LoadFiles
func LoadFiles(sourceFiles ...string) (err error) {
	return dc.LoadFiles(sourceFiles...)
}

// LoadExists
func LoadExists(sourceFiles ...string) (err error) {
	return dc.LoadExists(sourceFiles...)
}

// LoadData
func LoadData(dataSource ...interface{}) (err error) {
	return dc.LoadData(dataSource...)
}

// LoadSources
func LoadSources(format string, src []byte, more ...[]byte) (err error) {
	return dc.LoadSources(format, src, more...)
}

// LoadStrings
func LoadStrings(format string, str string, more ...string) (err error) {
	return dc.LoadStrings(format, str, more...)
}

/*************************************************************
 * helper functions for the default instance
 *************************************************************/

// Set
func Set(key string, val interface{}) (err error) {
	return dc.Set(key, val)
}

// WriteTo
func WriteTo(out io.Writer) (n int64, err error) {
	return dc.WriteTo(out)
}

// DumpTo
func DumpTo(out io.Writer, format string) (n int64, err error) {
	return dc.DumpTo(out, format)
}

// Data return all config data
func Data() map[string]interface{} {
	return dc.Data()
}

// ClearAll
func ClearAll()  {
	dc.ClearAll()
}

// ClearData
func ClearData()  {
	dc.ClearData()
}
