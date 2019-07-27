package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	at := assert.New(t)

	c := New("test")

	str := c.ToJSON()
	at.Equal("", str)

	err := c.LoadStrings(JSON, jsonStr)
	at.Nil(err)

	str = c.ToJSON()
	at.Contains(str, `"name":"app"`)

	buf := &bytes.Buffer{}
	_, err = c.WriteTo(buf)
	at.Nil(err)

	buf = &bytes.Buffer{}

	_, err = c.DumpTo(buf, "invalid")
	at.Error(err)

	_, err = c.DumpTo(buf, Yml)
	at.Error(err)

	_, err = c.DumpTo(buf, JSON)
	at.Nil(err)
}

