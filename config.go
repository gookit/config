package config

import (
	"sync"
)

// package version
const Version     = "1.0.1"

// supported config format
const (
	Json = "json"
	Yml  = "yml"
	Yaml = "yaml"
	Toml = "toml"
)

type stringArr []string
type stringMap map[string]string

// Decoder for decode yml,json,toml defFormat content
type Decoder func(blob []byte, v interface{}) (err error)
type Encoder func(v interface{}) (out []byte, err error)

// Options config options
type Options struct {
	// parse env value. like: "${EnvName}" "${EnvName|default}"
	ParseEnv bool
	// config is readonly
	Readonly bool
	// ignore key string case
	IgnoreCase bool
	// default format
	DefaultFormat string
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
	intCaches map[string]int
	strCaches map[string]string
	arrCaches map[string]stringArr
	mapCaches map[string]stringMap
}

// New
func New(name string) *Config {
	return &Config{
		name: name,
		data: make(map[string]interface{}),

		opts:  &Options{DefaultFormat: Json},
		encoders: map[string]Encoder{Json: JsonEncoder},
		decoders: map[string]Decoder{Json: JsonDecoder},
	}
}

// SetOptions
func (c *Config) SetOptions(opts *Options) {
	c.opts = opts

	if c.opts.DefaultFormat == "" {
		c.opts.DefaultFormat = Json
	}
}

// Readonly
func (c *Config) Readonly(readonly bool) {
	c.opts.Readonly = readonly
}

// Name get config name
func (c *Config) Name() string {
	return c.name
}

// Data get all config data
func (c *Config) Data() map[string]interface{} {
	return c.data
}

// SetDecoder
func (c *Config) SetDecoder(format string, decoder Decoder) {
	if format == Yml {
		format = Yaml
	}

	c.decoders[format] = decoder
}

// HasDecoder
func (c *Config) HasDecoder(format string) bool {
	_, ok := c.decoders[format]

	return ok
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
}

// ClearCaches
func (c *Config) ClearCaches() {
	c.intCaches = nil
	c.strCaches = nil
	c.mapCaches = nil
	c.arrCaches = nil
}

// initCaches
func (c *Config) initCaches() {
	c.intCaches = map[string]int{}
	c.strCaches = map[string]string{}
	c.arrCaches = map[string]stringArr{}
	c.mapCaches = map[string]stringMap{}
}
