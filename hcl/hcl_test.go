package hcl

import (
	"testing"

	"github.com/gookit/config/v2"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil/assert"
)

func TestDriver(t *testing.T) {
	is := assert.New(t)

	is.Eq("hcl", Driver.Name())

	c := config.NewEmpty("test")
	is.False(c.HasDecoder(config.Hcl))

	c.AddDriver(Driver)
	is.True(c.HasDecoder(config.Hcl))
	is.True(c.HasEncoder(config.Hcl))

	_, err := Encoder("some data")
	is.Err(err)
}

func TestLoadFile(t *testing.T) {
	is := assert.New(t)
	c := config.NewEmpty("test")
	c.AddDriver(Driver)

	err := c.LoadFiles("../testdata/hcl_base.hcl")
	is.NoErr(err)
	dump.Println(c.Data())

	err = c.LoadFiles("../testdata/hcl_example.conf")
	is.NoErr(err)

}
