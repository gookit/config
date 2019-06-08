package hcl

import (
	"testing"

	"github.com/gookit/config/v2"
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
