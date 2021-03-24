/*
Package config is a go config management implement. support YAML,TOML,JSON,INI,HCL format.

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/config

JSON format content example:

	{
		"name": "app",
		"debug": false,
		"baseKey": "value",
		"age": 123,
		"envKey": "${SHELL}",
		"envKey1": "${NotExist|defValue}",
		"map1": {
			"key": "val",
			"key1": "val1",
			"key2": "val2"
		},
		"arr1": [
			"val",
			"val1",
			"val2"
		],
		"lang": {
			"dir": "res/lang",
			"defLang": "en",
			"allowed": {
				"en": "val",
				"zh-CN": "val2"
			}
		}
	}

Usage please see example(more example please see examples folder in the lib):

*/
package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
)

// There are supported config format
const (
	Ini  = "ini"
	Hcl  = "hcl"
	Yml  = "yml"
	JSON = "json"
	Yaml = "yaml"
	Toml = "toml"

	// default delimiter
	defaultDelimiter byte = '.'
	// default struct tag name for binding data to struct
	defaultStructTag = "mapstructure"
)

// internal vars
type intArr []int
type strArr []string
type intMap map[string]int
type strMap map[string]string

// type fmtName string

// This is a default config manager instance
var dc = New("default")

// Driver interface
type Driver interface {
	Name() string
	GetDecoder() Decoder
	GetEncoder() Encoder
}

// Decoder for decode yml,json,toml format content
type Decoder func(blob []byte, v interface{}) (err error)

// Encoder for decode yml,json,toml format content
type Encoder func(v interface{}) (out []byte, err error)

// Options config options
type Options struct {
	// parse env value. like: "${EnvName}" "${EnvName|default}"
	ParseEnv bool
	// config is readonly
	Readonly bool
	// enable config data cache
	EnableCache bool
	// parse key, allow find value by key path. eg: 'key.sub' will find `map[key]sub`
	ParseKey bool
	// tag name for binding data to struct
	// Deprecated
	// please set tag name by DecoderConfig
	TagName string
	// the delimiter char for split key path, if `FindByPath=true`. default is '.'
	Delimiter byte
	// default write format
	DumpFormat string
	// default input format
	ReadFormat string
	// DecoderConfig setting for binding data to struct
	DecoderConfig *mapstructure.DecoderConfig
}

// Config structure definition
type Config struct {
	err error
	// config instance name
	name string
	lock sync.RWMutex

	// config options
	opts *Options
	// all config data
	data map[string]interface{}

	// loaded config files records
	loadedFiles  []string
	loadedDriver []string

	// decoders["toml"] = func(blob []byte, v interface{}) (err error){}
	// decoders["yaml"] = func(blob []byte, v interface{}) (err error){}
	// drivers map[string]Driver TODO Deprecated decoder and encoder, use driver instead
	decoders map[string]Decoder
	encoders map[string]Encoder

	// cache got config data
	intCache map[string]int
	strCache map[string]string

	iArrCache map[string]intArr
	iMapCache map[string]intMap
	sArrCache map[string]strArr
	sMapCache map[string]strMap
}

// New config instance
func New(name string) *Config {
	return &Config{
		name: name,
		opts: newDefaultOption(),
		data: make(map[string]interface{}),

		// default add JSON driver
		encoders: map[string]Encoder{JSON: JSONEncoder},
		decoders: map[string]Decoder{JSON: JSONDecoder},
	}
}

// NewEmpty config instance
func NewEmpty(name string) *Config {
	return &Config{
		name: name,
		opts: newDefaultOption(),
		data: make(map[string]interface{}),

		// don't add any drivers
		encoders: map[string]Encoder{},
		decoders: map[string]Decoder{},
	}
}

// NewWith create config instance, and you can call some init func
func NewWith(name string, fn func(c *Config)) *Config {
	return New(name).With(fn)
}

// NewWithOptions config instance
func NewWithOptions(name string, opts ...func(*Options)) *Config {
	return New(name).WithOptions(opts...)
}

// Default get the default instance
func Default() *Config {
	return dc
}

func newDefaultOption() *Options {
	return &Options{
		ParseKey:  true,
		TagName:   defaultStructTag,
		Delimiter: defaultDelimiter,
		// for export
		DumpFormat: JSON,
		ReadFormat: JSON,
		// struct decoder config
		DecoderConfig: newDefaultDecoderConfig(),
	}
}

func newDefaultDecoderConfig() *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		// tag name for binding struct
		TagName: defaultStructTag,
		// will auto convert string to int/uint
		WeaklyTypedInput: true,
	}
}

/*************************************************************
 * config setting
 *************************************************************/

// ParseEnv set parse env
func ParseEnv(opts *Options) { opts.ParseEnv = true }

// Readonly set readonly
func Readonly(opts *Options) { opts.Readonly = true }

// Delimiter set delimiter char
func Delimiter(sep byte) func(*Options) {
	return func(opts *Options) {
		opts.Delimiter = sep
	}
}

// EnableCache set readonly
func EnableCache(opts *Options) { opts.EnableCache = true }

// WithOptions with options
func WithOptions(opts ...func(*Options)) { dc.WithOptions(opts...) }

// WithOptions apply some options
func (c *Config) WithOptions(opts ...func(*Options)) *Config {
	if !c.IsEmpty() {
		panic("config: Cannot set options after data has been loaded")
	}

	// apply options
	for _, opt := range opts {
		opt(c.opts)
	}
	return c
}

// GetOptions get options
func GetOptions() *Options { return dc.Options() }

// Options get
func (c *Config) Options() *Options {
	return c.opts
}

// With apply some options
func (c *Config) With(fn func(c *Config)) *Config {
	fn(c)
	return c
}

// Readonly disable set data to config.
// Usage:
// 	config.LoadFiles(a, b, c)
// 	config.Readonly()
func (c *Config) Readonly() {
	c.opts.Readonly = true
}

/*************************************************************
 * config drivers
 *************************************************************/

// AddDriver set a decoder and encoder driver for a format.
func AddDriver(driver Driver) { dc.AddDriver(driver) }

// AddDriver set a decoder and encoder driver for a format.
func (c *Config) AddDriver(driver Driver) {
	format := driver.Name()
	c.decoders[format] = driver.GetDecoder()
	c.encoders[format] = driver.GetEncoder()
}

// HasDecoder has decoder
func (c *Config) HasDecoder(format string) bool {
	format = fixFormat(format)
	_, ok := c.decoders[format]
	return ok
}

// SetDecoder add/set a format decoder
// Deprecated
// please use driver instead
func SetDecoder(format string, decoder Decoder) {
	dc.SetDecoder(format, decoder)
}

// SetDecoder set decoder
// Deprecated
// please use driver instead
func (c *Config) SetDecoder(format string, decoder Decoder) {
	format = fixFormat(format)
	c.decoders[format] = decoder
}

// SetDecoders set decoders
// Deprecated
// please use driver instead
func (c *Config) SetDecoders(decoders map[string]Decoder) {
	for format, decoder := range decoders {
		c.SetDecoder(format, decoder)
	}
}

// SetEncoder set a encoder for the format
// Deprecated
// please use driver instead
func SetEncoder(format string, encoder Encoder) {
	dc.SetEncoder(format, encoder)
}

// SetEncoder set a encoder for the format
// Deprecated
// please use driver instead
func (c *Config) SetEncoder(format string, encoder Encoder) {
	format = fixFormat(format)
	c.encoders[format] = encoder
}

// SetEncoders set encoders
// Deprecated
// please use driver instead
func (c *Config) SetEncoders(encoders map[string]Encoder) {
	for format, encoder := range encoders {
		c.SetEncoder(format, encoder)
	}
}

// HasEncoder has encoder
func (c *Config) HasEncoder(format string) bool {
	format = fixFormat(format)
	_, ok := c.encoders[format]
	return ok
}

// DelDriver delete driver of the format
func (c *Config) DelDriver(format string) {
	format = fixFormat(format)

	if _, ok := c.decoders[format]; ok {
		delete(c.decoders, format)
	}

	if _, ok := c.encoders[format]; ok {
		delete(c.encoders, format)
	}
}

/*************************************************************
 * helper methods
 *************************************************************/

// Name get config name
func (c *Config) Name() string {
	return c.name
}

// Error get last error
func (c *Config) Error() error {
	return c.err
}

// IsEmpty of the config
func (c *Config) IsEmpty() bool {
	return len(c.data) == 0
}

// LoadedFiles get loaded files name
func (c *Config) LoadedFiles() []string {
	return c.loadedFiles
}

// ClearAll data and caches
func ClearAll() { dc.ClearAll() }

// ClearAll data and caches
func (c *Config) ClearAll() {
	c.ClearData()
	c.ClearCaches()

	c.loadedFiles = []string{}
	c.opts.Readonly = false
}

// ClearData clear data
func (c *Config) ClearData() {
	c.data = make(map[string]interface{})
	c.loadedFiles = []string{}
}

// ClearCaches clear caches
func (c *Config) ClearCaches() {
	if c.opts.EnableCache {
		c.intCache = nil
		c.strCache = nil
		c.sMapCache = nil
		c.sArrCache = nil
	}
}

/*************************************************************
 * helper methods/functions
 *************************************************************/

// record error
func (c *Config) addError(err error) {
	c.err = err
}

// format and record error
func (c *Config) addErrorf(format string, a ...interface{}) {
	c.err = fmt.Errorf(format, a...)
}

// GetEnv get os ENV value by name
// Deprecated
//	please use Getenv() instead
func GetEnv(name string, defVal ...string) (val string) {
	return Getenv(name, defVal...)
}

// Getenv get os ENV value by name. like os.Getenv, but support default value
// Notice:
// - Key is not case sensitive when getting
func Getenv(name string, defVal ...string) (val string) {
	if val = os.Getenv(name); val != "" {
		return
	}

	if len(defVal) > 0 {
		val = defVal[0]
	}
	return
}

// format key
func formatKey(key, sep string) string {
	return strings.Trim(strings.TrimSpace(key), sep)
}

// fix inc/conf/yaml format
func fixFormat(f string) string {
	if f == Yml {
		f = Yaml
	}

	if f == "inc" {
		f = Ini
	}

	// eg nginx config file.
	if f == "conf" {
		f = Hcl
	}

	return f
}
