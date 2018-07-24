package config

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var jsonStr = `{
    "name": "app",
    "debug": true,
    "baseKey": "value",
    "age": 123,
    "envKey": "${SHELL}",
    "envKey1": "${NotExist|defValue}",
    "map1": {
        "key": "val",
        "key1": "val1",
        "key2": "val2"
    },
    "arr1": [
        "val",
        "val1",
        "val2"
    ]
}`

func Example() {
	WithOptions(ParseEnv)

	// add Decoder and Encoder
	// use yaml github.com/gookit/config/yaml
	// AddDriver(Yaml, yaml.Driver)
	// use toml github.com/gookit/config/toml
	// AddDriver(Toml, toml.Driver)
	// use toml github.com/gookit/config/hcl
	// AddDriver(Hcl, hcl.Driver)
	// Or
	// config.DecoderEncoder(config.JSON, yaml.Decoder, yaml.Encoder)

	err := LoadFiles("testdata/json_base.json")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", Data())

	err = LoadFiles("testdata/json_other.json")
	// LoadFiles("testdata/json_base.json", "testdata/json_other.json")
	if err != nil {
		panic(err)
	}

	// load from string
	LoadSources(JSON, []byte(jsonStr))

	// fmt.Printf("config data: \n %#v\n", Data())
	fmt.Print("get config example:\n")

	name, ok := String("name")
	fmt.Printf("- get string\n ok: %v, val: %v\n", ok, name)

	arr1, ok := Strings("arr1")
	fmt.Printf("- get array\n ok: %v, val: %#v\n", ok, arr1)

	val0, ok := String("arr1.0")
	fmt.Printf("- get sub-value by path 'arr.index'\n ok: %v, val: %#v\n", ok, val0)

	map1, ok := StringMap("map1")
	fmt.Printf("- get map\n ok: %v, val: %#v\n", ok, map1)

	val0, ok = String("map1.key")
	fmt.Printf("- get sub-value by path 'map.key'\n ok: %v, val: %#v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", DefString("envKey1", ""))

	// set value
	Set("name", "new name")
	name, ok = String("name")
	fmt.Printf("- set string\n ok: %v, val: %v\n", ok, name)

	// if you want export config data
	// buf := new(bytes.Buffer)
	// _, err = config.DumpTo(buf, config.JSON)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("export config:\n%s", buf.String())

	// Out:
	// get config example:
	// - get string
	//  ok: true, val: app
	// - get array
	//  ok: true, val: []string{"val", "val1", "val2"}
	// - get sub-value by path 'arr.index'
	//  ok: true, val: "val"
	// - get map
	//  ok: true, val: map[string]string{"key":"val", "key1":"val1", "key2":"val2"}
	// - get sub-value by path 'map.key'
	//  ok: true, val: "val"
	// get env 'envKey' val: /bin/zsh
	// get env 'envKey1' val: defValue
	// - set string
	//  ok: true, val: new name
}

func ExampleConfig_DefBool() {
	// load from string
	LoadSources(JSON, []byte(jsonStr))

	val, ok := Bool("debug")
	fmt.Printf("get 'debug', ok: %v, val: %v\n", ok, val)
	val1 := DefBool("debug", false)
	fmt.Printf("get 'debug' with default, val: %v\n", val1)

	// Output:
	// get 'debug', ok: true, val: true
	// get 'debug' with default, val: true
}

func Example_exportConfig() {
	// Notice: before dump please set driver encoder
	// SetEncoder(Yaml, yaml.Encoder)

	ClearAll()
	// load from string
	LoadStrings(JSON, `{
"name": "app",
"age": 34
}`)

	buf := new(bytes.Buffer)
	_, err := DumpTo(buf, JSON)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", buf.String())

	// Output:
	// {"age":34,"name":"app"}
}

func BenchmarkGet(b *testing.B) {
	err := LoadStrings(JSON, jsonStr)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get("name")
	}
}

func TestBasic(t *testing.T) {
	st := assert.New(t)

	ClearAll()
	c := Default()
	st.True(c.HasDecoder(JSON))
	st.True(c.HasEncoder(JSON))
	st.Equal("default", c.Name())

	c = NewWithOptions("test", Readonly)
	opts := c.Options()
	st.True(opts.Readonly)
	st.Equal(JSON, opts.DumpFormat)
	st.Equal(JSON, opts.ReadFormat)
}

func TestLoad(t *testing.T) {
	st := assert.New(t)

	c := New("test")
	err := c.LoadExists("testdata/json_base.json", "not-exist.json")
	st.Nil(err)

	c.ClearAll()

	// load map
	err = c.LoadData(map[string]interface{}{
		"name":    "inhere",
		"age":     28,
		"working": true,
		"tags":    []string{"a", "b"},
		"info":    map[string]string{"k1": "a", "k2": "b"},
	})

	st.NotEmpty(c.Data())
	st.Nil(err)

	err = c.LoadData("invalid")
	st.Error(err)

	st.Panics(func() {
		c.WithOptions(ParseEnv)
	})

	err = c.LoadStrings(JSON, `{"name": "inhere"}`, jsonStr)
	st.Nil(err)

	err = c.LoadSources(JSON, []byte(`{"name": "inhere"}`), []byte(jsonStr))
	st.Nil(err)

	err = c.LoadSources(JSON, []byte(`invalid`))
	st.Error(err)

	err = c.LoadSources(JSON, []byte(`{"name": "inhere"}`), []byte(`invalid`))
	st.Error(err)

	c = New("test")

	err = c.LoadFiles("not-exist.json")
	st.Error(err)

	err = c.LoadFiles("testdata/json_error.json")
	st.Error(err)

	err = c.LoadExists("testdata/json_error.json")
	st.Error(err)

	err = c.LoadStrings("invalid", jsonStr)
	st.Error(err)
}

func TestJSONDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("json", JSONDriver.Name())

	// empty
	c := NewEmpty("test")
	st.False(c.HasDecoder(JSON))

	c.AddDriver(JSONDriver)
	st.True(c.HasDecoder(JSON))
	st.True(c.HasEncoder(JSON))
}

func TestDriver(t *testing.T) {
	st := assert.New(t)

	c := Default()
	st.True(c.HasDecoder(JSON))
	st.True(c.HasEncoder(JSON))

	c.DelDriver(JSON)
	st.False(c.HasDecoder(JSON))
	st.False(c.HasEncoder(JSON))

	AddDriver(JSONDriver)
	st.True(c.HasDecoder(JSON))
	st.True(c.HasEncoder(JSON))

	c.DelDriver(JSON)
	c.SetDecoders(map[string]Decoder{JSON: JSONDecoder})
	c.SetEncoders(map[string]Encoder{JSON: JSONEncoder})
	st.True(c.HasDecoder(JSON))
	st.True(c.HasEncoder(JSON))
}

func TestOptions(t *testing.T) {
	st := assert.New(t)

	// options: ParseEnv
	c := New("test")
	c.WithOptions(ParseEnv)

	st.True(c.Options().ParseEnv)

	err := c.LoadStrings(JSON, jsonStr)
	st.Nil(err)

	str, ok := c.String("envKey")
	st.True(ok)
	st.NotContains(str, "${")

	str, ok = c.String("envKey1")
	st.True(ok)
	st.NotContains(str, "${")
	st.Equal("defValue", str)

	// options: Readonly
	c = New("test")
	c.WithOptions(Readonly)

	st.True(c.Options().Readonly)

	err = c.LoadStrings(JSON, jsonStr)
	st.Nil(err)

	str, ok = c.String("name")
	st.True(ok)
	st.Equal("app", str)

	err = c.Set("name", "new app")
	st.Error(err)
}

func TestExport(t *testing.T) {
	at := assert.New(t)

	c := New("test")

	str := c.ToJSON()
	at.Equal("", str)

	c.LoadStrings(JSON, jsonStr)

	str = c.ToJSON()
	at.Contains(str, `"name":"app"`)

	buf := &bytes.Buffer{}
	_, err := c.WriteTo(buf)
	at.Nil(err)

	buf = &bytes.Buffer{}

	_, err = c.DumpTo(buf, "invalid")
	at.Error(err)

	_, err = c.DumpTo(buf, JSON)
	at.Nil(err)
}
