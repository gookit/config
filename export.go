package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

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
