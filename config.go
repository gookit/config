package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// Version is package version
const Version = "2.0.1"

// There are supported config format
const (
	Ini  = "ini"
	Hcl  = "hcl"
	Yml  = "yml"
	JSON = "json"
	Yaml = "yaml"
	Toml = "toml"
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
	// default write format
	DumpFormat string
	// default input format
	ReadFormat string
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
	loadedFiles []string

	// decoders["toml"] = func(blob []byte, v interface{}) (err error){}
	// decoders["yaml"] = func(blob []byte, v interface{}) (err error){}
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
		data: make(map[string]interface{}),

		// init options
		opts: &Options{DumpFormat: JSON, ReadFormat: JSON},

		// default add JSON driver
		encoders: map[string]Encoder{JSON: JSONEncoder},
		decoders: map[string]Decoder{JSON: JSONDecoder},
	}
}

// NewEmpty config instance
func NewEmpty(name string) *Config {
	return &Config{
		name: name,
		data: make(map[string]interface{}),

		// empty options
		opts: &Options{},

		// don't add any drivers
		encoders: map[string]Encoder{},
		decoders: map[string]Decoder{},
	}
}

// NewWithOptions config instance
func NewWithOptions(name string, opts ...func(*Options)) *Config {
	c := New(name)
	c.WithOptions(opts...)
	return c
}

// Default get the default instance
func Default() *Config {
	return dc
}

/*************************************************************
 * config setting
 *************************************************************/

// ParseEnv set parse env
func ParseEnv(opts *Options) {
	opts.ParseEnv = true
}

// Readonly set readonly
func Readonly(opts *Options) {
	opts.Readonly = true
}

// EnableCache set readonly
func EnableCache(opts *Options) {
	opts.EnableCache = true
}

// WithOptions with options
func WithOptions(opts ...func(*Options)) { dc.WithOptions(opts...) }

// WithOptions apply some options
func (c *Config) WithOptions(opts ...func(*Options)) {
	if !c.IsEmpty() {
		panic("config: Cannot set options after data has been loaded")
	}

	// apply options
	for _, opt := range opts {
		opt(c.opts)
	}
}

// GetOptions get options
func GetOptions() *Options { return dc.Options() }

// Options get
func (c *Config) Options() *Options {
	return c.opts
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
func SetDecoder(format string, decoder Decoder) {
	dc.SetDecoder(format, decoder)
}

// SetDecoder set decoder
func (c *Config) SetDecoder(format string, decoder Decoder) {
	format = fixFormat(format)
	c.decoders[format] = decoder
}

// SetDecoders set decoders
func (c *Config) SetDecoders(decoders map[string]Decoder) {
	for format, decoder := range decoders {
		c.SetDecoder(format, decoder)
	}
}

// SetEncoder set a encoder for the format
func SetEncoder(format string, encoder Encoder) {
	dc.SetEncoder(format, encoder)
}

// SetEncoder set a encoder for the format
func (c *Config) SetEncoder(format string, encoder Encoder) {
	format = fixFormat(format)
	c.encoders[format] = encoder
}

// SetEncoders set encoders
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

// Data return all config data
func Data() map[string]interface{} { return dc.Data() }

// Data get all config data
func (c *Config) Data() map[string]interface{} {
	return c.data
}

// Error get last error
func (c *Config) Error() error {
	return c.err
}

// IsEmpty of the config
func (c *Config) IsEmpty() bool {
	return len(c.data) == 0
}

// ToJSON string
func (c *Config) ToJSON() string {
	buf := &bytes.Buffer{}

	_, err := c.DumpTo(buf, JSON)
	if err != nil {
		return ""
	}

	return buf.String()
}

// WriteTo a writer
func WriteTo(out io.Writer) (int64, error) { return dc.WriteTo(out) }

// WriteTo Write out config data representing the current state to a writer.
func (c *Config) WriteTo(out io.Writer) (n int64, err error) {
	return c.DumpTo(out, c.opts.DumpFormat)
}

// DumpTo a writer and use format
func DumpTo(out io.Writer, format string) (int64, error) { return dc.DumpTo(out, format) }

// DumpTo use the format(json,yaml,toml) dump config data to a writer
func (c *Config) DumpTo(out io.Writer, format string) (n int64, err error) {
	var ok bool
	var encoder Encoder

	format = fixFormat(format)
	if encoder, ok = c.encoders[format]; !ok {
		err = errors.New("no exists or no register encoder for the format: " + format)
		return
	}

	// is empty
	if len(c.data) == 0 {
		return
	}

	// encode data to string
	encoded, err := encoder(&c.data)
	if err != nil {
		return
	}

	// write content to out
	num, err := fmt.Fprintln(out, string(encoded))
	if err != nil {
		return
	}

	return int64(num), nil
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
func GetEnv(name string, defVal ...string) (val string) {
	name = strings.ToUpper(name)
	if val = os.Getenv(name); val != "" {
		return
	}

	if len(defVal) > 0 {
		val = defVal[0]
	}
	return
}

// fix yaml format
func fixFormat(f string) string {
	if f == Yml {
		f = Yaml
	}

	return f
}
