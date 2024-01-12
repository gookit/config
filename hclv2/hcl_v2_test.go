package hclv2

import (
	"testing"

	"github.com/gookit/config/v2"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil/assert"
)

func TestDriver(t *testing.T) {
	is := assert.New(t)
	is.Eq("hcl", Driver.Name())
	is.Eq(config.Hcl, Driver.Name())

	c := config.NewEmpty("test")
	is.False(c.HasDecoder(config.Hcl))

	c.AddDriver(Driver)
	is.True(c.HasDecoder(config.Hcl))
	is.True(c.HasEncoder(config.Hcl))

	_, err := Encoder("some data")
	is.Err(err)
}

func TestHcl2Package(t *testing.T) {
	hclStr := `io_mode = "async"

service "http" "web_proxy" {
  listen_addr = "127.0.0.1:8080"
  
  process "main" {
    command = ["/usr/local/bin/awesome-app", "server"]
  }

  process "mgmt" {
    command = ["/usr/local/bin/awesome-app", "mgmt"]
  }
}`

	mp := make(map[string]any)
	// mp := make(map[string]cty.Type)
	// err := hclsimple.Decode("test.hcl", []byte(hclStr), nil, &mp)
	// assert.NoErr(t, err)
	dump.P(hclStr, mp)
	t.Skip("Not completed")
}

func TestLoadFile(t *testing.T) {
	is := assert.New(t)
	c := config.NewEmpty("test")
	c.AddDriver(Driver)
	is.True(c.HasDecoder(config.Hcl))

	t.Skip("Not completed")
	return
	err := c.LoadFiles("../testdata/hcl2_base.hcl")
	is.NoErr(err)
	dump.Println(c.Data())

	err = c.LoadFiles("../testdata/hcl2_example.hcl")
	is.NoErr(err)
}
