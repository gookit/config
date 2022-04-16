package config

import (
	"bytes"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	is := assert.New(t)

	c := New("test")

	str := c.ToJSON()
	is.Equal("", str)

	err := c.LoadStrings(JSON, jsonStr)
	is.Nil(err)

	str = c.ToJSON()
	is.Contains(str, `"name":"app"`)

	buf := &bytes.Buffer{}
	_, err = c.WriteTo(buf)
	is.Nil(err)

	buf = &bytes.Buffer{}

	_, err = c.DumpTo(buf, "invalid")
	is.Error(err)

	_, err = c.DumpTo(buf, Yml)
	is.Error(err)

	_, err = c.DumpTo(buf, JSON)
	is.Nil(err)
}

func TestConfig_Structure(t *testing.T) {
	is := assert.New(t)

	cfg := Default()
	cfg.ClearAll()

	err := cfg.LoadStrings(JSON, `{
"age": 28,
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)

	is.Nil(err)

	user := &struct {
		Age    int // always float64 from JSON
		Name   string
		Sports []string
	}{}
	// map all data
	err = MapStruct("", user)
	is.Nil(err)

	is.Equal(28, user.Age)
	is.Equal("inhere", user.Name)
	is.Equal("pingPong", user.Sports[0])

	// - auto convert string to int
	// age use string in JSON
	cfg1 := New("test")
	err = cfg1.LoadStrings(JSON, `{
"age": "26",
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)

	is.Nil(err)

	user1 := &struct {
		Age    int // always float64 from JSON
		Name   string
		Sports []string
	}{}
	err = cfg1.MapStruct("", user1)
	is.Nil(err)

	dump.P(*user1)

	// map some data
	err = cfg.LoadStrings(JSON, `{
"sec": {
	"key": "val",
	"age": 120,
	"tags": [12, 34]
}
}`)
	is.Nil(err)

	some := struct {
		Age  int
		Key  string
		Tags []int
	}{}
	err = BindStruct("sec", &some)
	is.Nil(err)
	is.Equal(120, some.Age)
	is.Equal(12, some.Tags[0])
	cfg.ClearAll()

	// custom data
	cfg = New("test")
	err = cfg.LoadData(map[string]interface{}{
		"key":  "val",
		"age":  120,
		"tags": []int{12, 34},
	})
	is.NoError(err)

	s1 := struct {
		Age  int
		Kye  string
		Tags []int
	}{}
	err = cfg.BindStruct("", &s1)
	is.Nil(err)
	is.Equal(120, s1.Age)
	is.Equal(12, s1.Tags[0])

	// key not exist
	err = cfg.BindStruct("not-exist", &s1)
	is.Error(err)
	is.Equal("this key does not exist in the config", err.Error())

	// invalid dst
	err = cfg.BindStruct("sec", "invalid")
	is.Error(err)

	cfg.ClearAll()
}

func TestMapStruct_embedded_struct_squash_false(t *testing.T) {
	loader := NewWithOptions("test", func(options *Options) {
		options.DecoderConfig.TagName = "json"
		options.DecoderConfig.Squash = false
	})
	assert.False(t, loader.Options().DecoderConfig.Squash)

	err := loader.LoadStrings(JSON, `
{
  "c": "12",
  "test1": {
	"b": "34"
  }
}
`)
	assert.NoError(t, err)
	dump.Println(loader.Data())
	assert.Equal(t, 12, loader.Int("c"))
	assert.Equal(t, 34, loader.Int("test1.b"))

	type Test1 struct {
		B int `json:"b"`
	}
	type Test2 struct {
		Test1
		C int `json:"c"`
	}
	cfg := &Test2{}

	err = loader.MapStruct("", cfg)
	assert.NoError(t, err)
	dump.Println(cfg)
	assert.Equal(t, 34, cfg.Test1.B)

	type Test3 struct {
		*Test1
		C int `json:"c"`
	}
	cfg1 := &Test3{}
	err = loader.MapStruct("", cfg1)
	assert.NoError(t, err)
	dump.Println(cfg1)
	assert.Equal(t, 34, cfg1.Test1.B)
}

func TestMapStruct_embedded_struct_squash_true(t *testing.T) {
	loader := NewWithOptions("test", func(options *Options) {
		options.DecoderConfig.TagName = "json"
		options.DecoderConfig.Squash = true
	})
	assert.True(t, loader.Options().DecoderConfig.Squash)

	err := loader.LoadStrings(JSON, `
{
  "c": "12",
  "test1": {
	"b": "34"
  }
}
`)
	assert.NoError(t, err)
	dump.Println(loader.Data())
	assert.Equal(t, 12, loader.Int("c"))
	assert.Equal(t, 34, loader.Int("test1.b"))

	type Test1 struct {
		B int `json:"b"`
	}
	type Test2 struct {
		Test1
		// Test1 `json:",squash"`
		C int `json:"c"`
	}
	cfg := &Test2{}

	err = loader.MapStruct("", cfg)
	assert.NoError(t, err)
	dump.Println(cfg)
	assert.Equal(t, 0, cfg.Test1.B)

	type Test3 struct {
		*Test1
		C int `json:"c"`
	}
	cfg1 := &Test3{}
	err = loader.MapStruct("", cfg1)
	assert.NoError(t, err)
	dump.Println(cfg1)
	assert.Equal(t, 34, cfg1.Test1.B)
}

func TestMapOnExists(t *testing.T) {
	cfg := Default()
	cfg.ClearAll()

	err := cfg.LoadStrings(JSON, `{
"age": 28,
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)
	assert.NoError(t, err)
	assert.NoError(t, MapOnExists("not-exists", nil))

	user := &struct {
		Age    int
		Name   string
		Sports []string
	}{}
	assert.NoError(t, MapOnExists("", user))

	assert.Equal(t, 28, user.Age)
	assert.Equal(t, "inhere", user.Name)
}

func TestConfig_BindStruct_set_DecoderConfig(t *testing.T) {
	cfg := NewWith("test", func(c *Config) {
		c.opts.DecoderConfig = nil
	})
	err := cfg.LoadStrings(JSON, `{
"age": 28,
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)
	assert.NoError(t, err)

	user := &struct {
		Age    int
		Name   string
		Sports []string
	}{}
	assert.NoError(t, cfg.BindStruct("", user))

	assert.Equal(t, 28, user.Age)
	assert.Equal(t, "inhere", user.Name)

	// not use ptr
	assert.Error(t, cfg.BindStruct("", *user))
}

func TestConfig_BindStruct_error(t *testing.T) {
	// cfg := NewEmpty()
}
