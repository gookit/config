package config

import (
	"sync"
)

// const Ini = "ini"
const Json = "json"
const Yml = "yml"
const Yaml = "yaml"
const Toml = "toml"

const Version = "0.0.1"
const DefaultNode = "__DEFAULT"

// Decoder for decode yml,json,toml defFormat content
type Decoder func(blob []byte, v interface{}) (err error)
type Encoder func(v interface{}) (out []byte, err error)

// Node
type Node struct {
	name string

	isArray   bool
	mapValues map[string]string
	arrValues map[string][]string

	// sub nodes
	hasChild bool
	nodes    map[string]*Node
}

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
	name string
	lock sync.RWMutex

	data  map[string]interface{}
	nodes map[string]*Node

	options   *Options
	readOnly  bool
	// only load exists file
	loadExist bool
	// default format
	defFormat string
	// ignore key string case
	ignoreCase bool

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

func (c *Config) IgnoreCase(ignoreCase bool) {
	c.ignoreCase = ignoreCase
}

func (c *Config) ReadOnly(readOnly bool) {
	c.readOnly = readOnly
}

func (c *Config) Data() map[string]interface{} {
	return c.data
}

func (c *Config) HasDecoder(format string) bool {
	_, ok := c.decoders[format]

	return ok
}
