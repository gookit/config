package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/mitchellh/mapstructure"
)

// MapStruct alias method of the 'Structure'
func MapStruct(key string, dst interface{}) error { return dc.MapStruct(key, dst) }

// MapStruct alias method of the 'Structure'
func (c *Config) MapStruct(key string, dst interface{}) error {
	return c.Structure(key, dst)
}

// BindStruct alias method of the 'Structure'
func BindStruct(key string, dst interface{}) error { return dc.BindStruct(key, dst) }

// BindStruct alias method of the 'Structure'
func (c *Config) BindStruct(key string, dst interface{}) error {
	return c.Structure(key, dst)
}

// Structure get config data and binding to the dst structure.
// Usage:
// 	dbInfo := Db{}
// 	config.Structure("db", &dbInfo)
func (c *Config) Structure(key string, dst interface{}) (err error) {
	var data interface{}

	// binding all data
	if key == "" {
		data = c.data
	} else { // some data of the config
		var ok bool

		data, ok = c.GetValue(key)
		if !ok {
			return errNotFound
		}
	}

	// format := JSON
	// if len(driverName) > 0 {
	// 	format = driverName[0]
	// }

	// decoder := c.getDecoderByFormat(format)
	// if decoder != nil {
	//
	// } else { // try use mergo
	// 	return mergo.Map(dst, data)
	// }

	err = mapstructure.Decode(data, dst)
	return
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
