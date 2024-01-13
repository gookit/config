package config

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/jsonutil"
	"github.com/gookit/goutil/testutil/assert"
)

func TestExport(t *testing.T) {
	is := assert.New(t)
	c := New("test")

	str := c.ToJSON()
	is.Eq("", str)

	err := c.LoadStrings(JSON, jsonStr)
	is.Nil(err)

	str = c.ToJSON()
	is.Contains(str, `"name":"app"`)

	buf := &bytes.Buffer{}
	_, err = c.WriteTo(buf)
	is.Nil(err)

	// test dump
	buf = &bytes.Buffer{}
	_, err = c.DumpTo(buf, "invalid")
	is.Err(err)
	_, err = c.DumpTo(buf, Yml)
	is.Err(err)

	_, err = c.DumpTo(buf, JSON)
	is.Nil(err)
}

func TestDumpTo_encode_error(t *testing.T) {
	is := assert.New(t)
	c := NewEmpty("test")
	is.NoErr(c.Set("age", 34))

	drv := NewDriver(JSON, JSONDecoder, func(v any) (out []byte, err error) {
		return nil, errors.New("encode data error")
	})
	c.WithDriver(drv)

	// encode error
	buf := &bytes.Buffer{}
	_, err := c.DumpTo(buf, JSON)
	is.ErrMsg(err, "encode data error")

	is.Empty(c.ToJSON())
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
	type User struct {
		Age    int // always float64 from JSON
		Name   string
		Sports []string
	}

	user := &User{}
	// map all data
	err = MapStruct("", user)
	is.Nil(err)

	is.Eq(28, user.Age)
	is.Eq("inhere", user.Name)
	is.Eq("pingPong", user.Sports[0])

	// map all data
	u1 := &User{}
	err = Decode(u1)
	is.Nil(err)

	is.Eq(28, u1.Age)
	is.Eq("inhere", u1.Name)
	is.Eq("pingPong", u1.Sports[0])

	// - auto convert string to int
	// age use string in JSON
	cfg1 := New("test")
	err = cfg1.LoadStrings(JSON, `{
"age": "26",
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)

	is.Nil(err)

	user1 := &User{}
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
	is.Eq(120, some.Age)
	is.Eq(12, some.Tags[0])
	cfg.ClearAll()

	// custom data
	cfg = New("test")
	err = cfg.LoadData(map[string]any{
		"key":  "val",
		"age":  120,
		"tags": []int{12, 34},
	})
	is.NoErr(err)

	s1 := struct {
		Age  int
		Kye  string
		Tags []int
	}{}
	err = cfg.BindStruct("", &s1)
	is.Nil(err)
	is.Eq(120, s1.Age)
	is.Eq(12, s1.Tags[0])

	// key not exist
	err = cfg.BindStruct("not-exist", &s1)
	is.Err(err)
	is.Eq("this key does not exist in the config", err.Error())

	// invalid dst
	err = cfg.BindStruct("sec", "invalid")
	is.Err(err)

	cfg.ClearAll()
}

func TestMapStruct_embedded_struct_squash_false(t *testing.T) {
	loader := NewWithOptions("test", func(options *Options) {
		options.DecoderConfig.TagName = "json"
		options.DecoderConfig.Squash = false
	})
	assert.False(t, loader.Options().DecoderConfig.Squash)

	err := loader.LoadStrings(JSON, `{
  "c": "12",
  "test1": {
	"b": "34"
  }
}`)
	assert.NoErr(t, err)
	dump.Println(loader.Data())
	assert.Eq(t, 12, loader.Int("c"))
	assert.Eq(t, 34, loader.Int("test1.b"))

	type Test1 struct {
		B int `json:"b"`
	}
	type Test2 struct {
		Test1
		C int `json:"c"`
	}
	cfg := &Test2{}

	err = loader.MapStruct("", cfg)
	assert.NoErr(t, err)
	dump.Println(cfg)
	assert.Eq(t, 34, cfg.Test1.B)

	type Test3 struct {
		*Test1
		C int `json:"c"`
	}
	cfg1 := &Test3{}
	err = loader.MapStruct("", cfg1)
	assert.NoErr(t, err)
	dump.Println(cfg1)
	assert.Eq(t, 34, cfg1.Test1.B)

	loader.SetData(map[string]any{
		"c": 120,
		"b": 340,
	})
	dump.Println(loader.Data())

	cfg2 := &Test3{}
	err = loader.BindStruct("", cfg2)

	cfg3 := &Test3{}
	_ = jsonutil.DecodeString(`{"c": 12, "b": 34}`, cfg3)

	dump.Println(cfg2, cfg3)
}

func TestMapStruct_embedded_struct_squash_true(t *testing.T) {
	loader := NewWithOptions("test", func(options *Options) {
		options.DecoderConfig.TagName = "json"
		options.DecoderConfig.Squash = true
	})
	assert.True(t, loader.Options().DecoderConfig.Squash)

	err := loader.LoadStrings(JSON, `{
  "c": "12",
  "test1": {
	"b": "34"
  }
}`)
	assert.NoErr(t, err)
	dump.Println(loader.Data())
	assert.Eq(t, 12, loader.Int("c"))
	assert.Eq(t, 34, loader.Int("test1.b"))

	type Test1 struct {
		B int `json:"b"`
	}

	// use value - will not set ok
	type Test2 struct {
		Test1
		// Test1 `json:",squash"`
		C int `json:"c"`
	}
	cfg := &Test2{}

	err = loader.MapStruct("", cfg)
	assert.NoErr(t, err)
	dump.Println(cfg)
	assert.Eq(t, 0, cfg.Test1.B)

	// use pointer
	type Test3 struct {
		*Test1
		C int `json:"c"`
	}
	cfg1 := &Test3{}
	err = loader.MapStruct("", cfg1)
	assert.NoErr(t, err)
	dump.Println(cfg1)
	assert.Eq(t, 34, cfg1.B)
	assert.Eq(t, 34, cfg1.Test1.B)

	loader.SetData(map[string]any{
		"c": 120,
		"b": 340,
	})
	dump.Println(loader.Data())

	cfg2 := &Test3{}
	err = loader.BindStruct("", cfg2)

	cfg3 := &Test3{}
	_ = jsonutil.DecodeString(`{"c": 12, "b": 34}`, cfg3)

	dump.Println(cfg2, cfg3)
}

func TestMapOnExists(t *testing.T) {
	cfg := Default()
	cfg.ClearAll()

	err := cfg.LoadStrings(JSON, `{
"age": 28,
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)
	assert.NoErr(t, err)
	assert.NoErr(t, MapOnExists("not-exists", nil))

	user := &struct {
		Age    int
		Name   string
		Sports []string
	}{}
	assert.NoErr(t, MapOnExists("", user))

	assert.Eq(t, 28, user.Age)
	assert.Eq(t, "inhere", user.Name)
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
	assert.NoErr(t, err)

	user := &struct {
		Age    int
		Name   string
		Sports []string
	}{}
	assert.NoErr(t, cfg.BindStruct("", user))

	assert.Eq(t, 28, user.Age)
	assert.Eq(t, "inhere", user.Name)

	// not use ptr
	assert.Err(t, cfg.BindStruct("", *user))
}

func TestConfig_BindStruct_error(t *testing.T) {
	// cfg := NewEmpty()
}

func TestConfig_BindStruct_default(t *testing.T) {
	type MyConf struct {
		Env   string `default:"${APP_ENV | dev}"`
		Debug bool   `default:"${APP_DEBUG | false}"`
	}

	cfg := NewWithOptions("test", ParseEnv, ParseDefault)
	// cfg.SetData(map[string]any{
	// 	"env": "prod",
	// 	"debug": "true",
	// })

	mc := &MyConf{}
	err := cfg.Decode(mc)
	dump.P(mc)
	assert.NoErr(t, err)
	assert.Eq(t, "dev", mc.Env)
	assert.False(t, mc.Debug)
}
