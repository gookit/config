package config

import (
	"sync"
	"io"
	"fmt"
	"errors"
)

// package version
const Version = "1.0.3"

// supported config format
const (
	Ini = "ini"
	Hcl  = "hcl"
	Yml  = "yml"
	Json = "json"
	Yaml = "yaml"
	Toml = "toml"
)

// internal vars
type intArr []int
type strArr []string
type intMap map[string]int
type strMap map[string]string

// Driver
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

// Config
type Config struct {
	// config instance name
	name string
	lock sync.RWMutex

	// config options
	opts *Options
	// all config data
	data map[string]interface{}

	// loaded config files
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

// New
func New(name string) *Config {
	return &Config{
		name: name,
		data: make(map[string]interface{}),

		// init options
		opts: &Options{DumpFormat: Json, ReadFormat: Json},

		// default add json driver
		encoders: map[string]Encoder{Json: JsonEncoder},
		decoders: map[string]Decoder{Json: JsonDecoder},
	}
}

/*************************************************************
 * config setting
 *************************************************************/

// SetOptions
func (c *Config) SetOptions(opts *Options) {
	c.opts = opts

	if c.opts.DumpFormat == "" {
		c.opts.DumpFormat = Json
	}

	if c.opts.ReadFormat == "" {
		c.opts.ReadFormat = Json
	}
}

// Name get config name
func (c *Config) Name() string {
	return c.name
}

// Data get all config data
func (c *Config) Data() map[string]interface{} {
	return c.data
}

// Readonly
func (c *Config) Readonly(readonly bool) {
	c.opts.Readonly = readonly
}

// AddDriver set a decoder and encoder driver for a format.
func (c *Config) AddDriver(format string, driver Driver) {
	format = fixFormat(format)
	if format != driver.Name() {
		panic(fmt.Sprintf(
			"format name must be equals to the driver name. current format:%s driver:%s",
			format,
			driver.Name(),
		))
	}

	c.decoders[format] = driver.GetDecoder()
	c.encoders[format] = driver.GetEncoder()
}

// DecoderEncoder set a decoder and encoder for a format.
func (c *Config) DecoderEncoder(format string, decoder Decoder, encoder Encoder) {
	format = fixFormat(format)

	c.decoders[format] = decoder
	c.encoders[format] = encoder
}

// HasDecoder
func (c *Config) HasDecoder(format string) bool {
	format = fixFormat(format)
	_, ok := c.decoders[format]

	return ok
}

// SetDecoder
func (c *Config) SetDecoder(format string, decoder Decoder) {
	format = fixFormat(format)
	c.decoders[format] = decoder
}

// SetDecoders
func (c *Config) SetDecoders(decoders map[string]Decoder) {
	for format, decoder := range decoders {
		c.SetDecoder(format, decoder)
	}
}

// SetEncoder
func (c *Config) SetEncoder(format string, encoder Encoder) {
	format = fixFormat(format)
	c.encoders[format] = encoder
}

// SetEncoders
func (c *Config) SetEncoders(encoders map[string]Encoder) {
	for format, encoder := range encoders {
		c.SetEncoder(format, encoder)
	}
}

// HasEncoder
func (c *Config) HasEncoder(format string) bool {
	format = fixFormat(format)
	_, ok := c.encoders[format]

	return ok
}

/*************************************************************
 * helper methods
 *************************************************************/

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

// ClearAll
func (c *Config) ClearAll() {
	c.ClearData()
	c.ClearCaches()

	c.loadedFiles = []string{}
}

// ClearData
func (c *Config) ClearData() {
	c.data = make(map[string]interface{})
	c.loadedFiles = []string{}
}

// ClearCaches
func (c *Config) ClearCaches() {
	c.intCache = nil
	c.strCache = nil
	c.sMapCache = nil
	c.sArrCache = nil
}

// initCaches
func (c *Config) initCaches() {
	c.intCache = map[string]int{}
	c.strCache = map[string]string{}
	c.sArrCache = map[string]strArr{}
	c.sMapCache = map[string]strMap{}
}

// fixFormat
func fixFormat(f string) string {
	if f == Yml {
		f = Yaml
	}

	return f
}
