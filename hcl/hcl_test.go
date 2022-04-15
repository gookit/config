package hcl

import (
	"testing"

	"github.com/gookit/config/v2"
	"github.com/gookit/goutil/dump"
	"github.com/stretchr/testify/assert"
)

func TestDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("hcl", Driver.Name())

	c := config.NewEmpty("test")
	st.False(c.HasDecoder(config.Hcl))

	c.AddDriver(Driver)
	st.True(c.HasDecoder(config.Hcl))
	st.True(c.HasEncoder(config.Hcl))

	_, err := Encoder("some data")
	st.Error(err)
}

func TestLoadFile(t *testing.T) {
	is := assert.New(t)
	c := config.NewEmpty("test")
	c.AddDriver(Driver)

	err := c.LoadFiles("../testdata/hcl_base.hcl")
	is.NoError(err)
	dump.Println(c.Data())

	err = c.LoadFiles("../testdata/hcl_example.conf")
	is.NoError(err)

}
