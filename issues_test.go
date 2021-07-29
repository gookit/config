package config_test

import (
	"testing"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil"
	"github.com/stretchr/testify/assert"
)

// https://github.com/gookit/config/issues/37
func TestIssues_37(t *testing.T) {
	is := assert.New(t)

	c := config.New("test")
	c.AddDriver(yaml.Driver)

	err := c.LoadStrings(config.JSON, `
{
    "lang": {
        "allowed": {
            "en": "ddd"
        }
    }
}
`)
	is.NoError(err)
	dump.Println(c.Data())

	is.Panics(func() {
		_ = c.LoadStrings(config.Yaml, `
lang:
  allowed:
    en: "666"
`)
	})
}

// https://github.com/gookit/config/issues/37
func TestIssues37_yaml_v3(t *testing.T) {
	is := assert.New(t)

	c := config.New("test")
	c.AddDriver(yamlv3.Driver)

	err := c.LoadStrings(config.JSON, `
{
    "lang": {
        "allowed": {
            "en": "ddd"
        }
    }
}
`)
	is.NoError(err)
	dump.Println(c.Data())

	err = c.LoadStrings(config.Yaml, `
lang:
  newKey: hhh
  allowed:
    en: "666"
`)
	is.NoError(err)
	dump.Println(c.Data())
}

// BindStruct doesn't seem to work with env var substitution
// https://github.com/gookit/config/issues/46
func TestIssues_46(t *testing.T) {
	is := assert.New(t)

	c := config.New("test").WithOptions(config.ParseEnv)
	err := c.LoadStrings(config.JSON, `
{
  "http": {
    "port": "${HTTP_PORT|8080}"
  }
}
`)

	is.NoError(err)
	dump.Println(c.Data())

	val, _ := c.GetValue("http")
	mp := val.(map[string]interface{})
	dump.Println(mp)
	is.Equal("${HTTP_PORT|8080}", mp["port"])

	smp := c.StringMap("http")
	dump.Println(smp)
	is.Contains(smp, "port")
	is.Equal("8080", smp["port"])

	type Http struct {
		Port int
	}

	h := &Http{}
	err = c.BindStruct("http", h)
	is.NoError(err)
	dump.Println(h)
	is.Equal(8080, h.Port)

	testutil.MockEnvValue("HTTP_PORT", "19090", func(_ string) {
		h := &Http{}
		err = c.BindStruct("http", h)
		is.NoError(err)
		dump.Println(h)
		is.Equal(19090, h.Port)
	})

}
