package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
)

// Version is package version
const Version = "1.0.10"

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
	// config instance name
	name string
	lock sync.RWMutex

	// config options
	opts *Options
	// all config data
	data map[string]interface{}

	// loaded config files records
	loadedFiles []string
	initialized bool

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

/*************************************************************
 * config setting
 *************************************************************/

// Options get
func (c *Config) Options() *Options {
	return c.opts
}

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

// WithOptions apply some options
func (c *Config) WithOptions(opts ...func(*Options)) {
	if c.initialized {
		panic("config: Cannot set options after initialization is complete")
	}

	// apply options
	for _, opt := range opts {
		opt(c.opts)
	}
}

/*************************************************************
 * config drivers
 *************************************************************/

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

// Data get all config data
func (c *Config) Data() map[string]interface{} {
	return c.data
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

// WriteTo Write out config data representing the current state to a writer.
func (c *Config) WriteTo(out io.Writer) (n int64, err error) {
	return c.DumpTo(out, c.opts.DumpFormat)
}

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
func (c *Config) ClearAll() {
	c.ClearData()
	c.ClearCaches()

	c.loadedFiles = []string{}
}

// ClearData clear data
func (c *Config) ClearData() {
	c.data = make(map[string]interface{})
	c.loadedFiles = []string{}
}

// ClearCaches clear caches
func (c *Config) ClearCaches() {
	c.intCache = nil
	c.strCache = nil
	c.sMapCache = nil
	c.sArrCache = nil
}

// fixFormat
func fixFormat(f string) string {
	if f == Yml {
		f = Yaml
	}

	return f
}
