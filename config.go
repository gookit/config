package config

import (
	"sync"
)

// supported config format
const (
	Json = "json"
	Yml  = "yml"
	Yaml = "yaml"
	Toml = "toml"
)

const (
	Version     = "0.0.1"
	DefaultNode = "__DEFAULT"
)

type stringMap map[string]string
type stringArr []string

// Decoder for decode yml,json,toml defFormat content
type Decoder func(blob []byte, v interface{}) (err error)
type Encoder func(v interface{}) (out []byte, err error)

// Options
type Options struct {
	// config is readonly
	Readonly bool
	// ignore key string case
	IgnoreCase bool
	// default format
	DefaultFormat string
	// only load exists file
	IgnoreNotExist bool
}

// Config
type Config struct {
	// config instance name
	name string
	lock sync.RWMutex

	// all config data
	data map[string]interface{}

	// cache got config data
	intCaches map[string]int
	strCaches map[string]string
	arrCaches map[string]stringArr
	mapCaches map[string]stringMap

	options  *Options
	readOnly bool
	// only load exists file
	loadExist bool
	// default format
	defFormat string
	// ignore key string case
	// ignoreCase bool

	loadedFiles []string

	// decoders["toml"] = func(data string, v interface{}) (err error){}
	// decoders["yaml"] = func(data string, v interface{}) (err error){}
	decoders map[string]Decoder
	encoders map[string]Encoder
}

// New
func New(name string) *Config {
	return &Config{
		name: name,
		data: make(map[string]interface{}),

		defFormat: Json,

		encoders: map[string]Encoder{Json: JsonEncoder},
		decoders: map[string]Decoder{Json: JsonDecoder},
	}
}

// SetDecoder
func (c *Config) SetDecoder(format string, decoder Decoder) {
	c.decoders[format] = decoder
}

// ReadOnly
func (c *Config) ReadOnly(readOnly bool) {
	c.readOnly = readOnly
}

// Name get config name
func (c *Config) Name() string {
	return c.name
}

// Data get all config data
func (c *Config) Data() map[string]interface{} {
	return c.data
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
