package hcl

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/gookit/config"
)

func TestDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("hcl", Driver.Name())

	c := config.NewEmpty("test")
	st.False(c.HasDecoder(config.Hcl))
	st.Panics(func() {
		c.AddDriver("invalid", Driver)
	})
	c.AddDriver(config.Hcl, Driver)
	st.True(c.HasDecoder(config.Hcl))
	st.True(c.HasEncoder(config.Hcl))
}

