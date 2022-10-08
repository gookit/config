package hclv2

import (
	"testing"

	"github.com/gookit/config/v2"
	"github.com/gookit/goutil/dump"
	"github.com/stretchr/testify/assert"
)

func TestDriver(t *testing.T) {
	is := assert.New(t)
	is.Equal("hcl", Driver.Name())
	is.Equal(config.Hcl, Driver.Name())

	c := config.NewEmpty("test")
	is.False(c.HasDecoder(config.Hcl))

	c.AddDriver(Driver)
	is.True(c.HasDecoder(config.Hcl))
	is.True(c.HasEncoder(config.Hcl))

	_, err := Encoder("some data")
	is.Error(err)
}

func TestLoadFile(t *testing.T) {
	is := assert.New(t)
	c := config.NewEmpty("test")
	c.AddDriver(Driver)
	is.True(c.HasDecoder(config.Hcl))

	t.Skip("Not completed")
	return
	err := c.LoadFiles("../testdata/hcl2_base.hcl")
	is.NoError(err)
	dump.Println(c.Data())

	err = c.LoadFiles("../testdata/hcl2_example.hcl")
	is.NoError(err)
}
