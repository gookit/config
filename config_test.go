package config

import (
	"fmt"
	"bytes"
	"testing"
	"github.com/stretchr/testify/assert"
)

var jsonStr = `
{
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
}
`

func Example() {
	SetOptions(&Options{
		ParseEnv: true,
	})

	// add Decoder and Encoder
	// use yaml github.com/gookit/config/yaml
	// AddDriver(Yaml, yaml.Driver)
	// use toml github.com/gookit/config/toml
	// AddDriver(Toml, toml.Driver)
	// use toml github.com/gookit/config/hcl
	// AddDriver(Hcl, hcl.Driver)
	// Or
	// config.DecoderEncoder(config.Json, yaml.Decoder, yaml.Encoder)

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
	LoadSources(Json, []byte(jsonStr))

	// fmt.Printf("config data: \n %#v\n", Data())
	fmt.Print("get config example:\n")

	name, ok := GetString("name")
	fmt.Printf("- get string\n ok: %v, val: %v\n", ok, name)

	arr1, ok := GetStringArr("arr1")
	fmt.Printf("- get array\n ok: %v, val: %#v\n", ok, arr1)

	val0, ok := GetString("arr1.0")
	fmt.Printf("- get sub-value by path 'arr.index'\n ok: %v, val: %#v\n", ok, val0)

	map1, ok := GetStringMap("map1")
	fmt.Printf("- get map\n ok: %v, val: %#v\n", ok, map1)

	val0, ok = GetString("map1.key")
	fmt.Printf("- get sub-value by path 'map.key'\n ok: %v, val: %#v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", DefString("envKey1", ""))

	// set value
	Set("name", "new name")
	name, ok = GetString("name")
	fmt.Printf("- set string\n ok: %v, val: %v\n", ok, name)

	// if you want export config data
	// buf := new(bytes.Buffer)
	// _, err = config.DumpTo(buf, config.Json)
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
	LoadSources(Json, []byte(jsonStr))

	val, ok := GetBool("debug")
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
	LoadStrings(Json, `{
"name": "app",
"age": 34
}`)

	buf := new(bytes.Buffer)
	_, err := DumpTo(buf, Json)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", buf.String())

	// Output:
	// {"age":34,"name":"app"}
}

func BenchmarkGet(b *testing.B) {
	err := LoadStrings(Json, jsonStr)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get("name")
	}
}

func TestLoadData(t *testing.T) {
	st := assert.New(t)

	ClearAll()

	// load map
	err := LoadData(map[string]interface{}{
		"name": "inhere",
		"age": 28,
		"working": true,
		"tags": []string{"a", "b"},
		"info": map[string]string{"k1":"a", "k2":"b"},
	})

	st.NotEmpty(Data())

	if st.Nil(err) {
		str, ok := GetString("name")
		st.True(ok)
		st.Equal("inhere", str)

		str, ok = GetString("notExists")
		st.False(ok)
		st.Equal("", str)

		def := DefString("notExists", "defVal")
		st.Equal("defVal", def)
	}
}
