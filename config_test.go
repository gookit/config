package config

import (
	"fmt"
	"github.com/gookit/config/yaml"
	"bytes"
	"github.com/gookit/config/toml"
)

var yamlStr = `
name: app2
debug: false
age: 23
baseKey: value2

map1:
    key: val2
    key2: val20

arr1:
    - val1
    - val21
`

func Example_useYaml() {
	// add yaml decoder
	SetDecoder(Yaml, yaml.Decoder)
	err := LoadFiles("testdata/yml_other.yml")
	if err != nil {
		panic(err)
	}

	// load from string
	LoadSources(Yaml, []byte(yamlStr))

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


	// Output:
	// get config example:
	// - get string
	// ok: true, val: app2
	// - get array
	// ok: true, val: []string{"val1", "val21"}
	// - get sub-value by path 'arr.index'
	// ok: true, val: "val1"
	// - get map
	// ok: true, val: map[string]string{"key":"val2", "key2":"val20"}
	// - get sub-value by path 'map.key'
	// ok: true, val: "val2"
	// get env 'envKey' val: /bin/zsh
	// get env 'envKey1' val: defValue
}

var tomlStr = `
title = "TOML Example"
name = "app"

envKey = "${SHELL}"
envKey1 = "${NotExist|defValue}"

arr1 = [
  "alpha",
  "omega"
]

[map1]
name = "inhere"
org = "GitHub"
`

func Example_useToml() {
	SetOptions(&Options{
		ParseEnv: true,
	})
	SetDriver(Toml, toml.Decoder, toml.Encoder)

	err := LoadFiles("testdata/toml_base.toml")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", Data())

	// load more files
	err = LoadFiles("testdata/toml_other.toml")
	// can also
	// LoadFiles("testdata/toml_base.toml", "testdata/toml_other.toml")
	if err != nil {
		panic(err)
	}

	// load from string
	LoadSources(Toml, []byte(tomlStr))

	// fmt.Printf("config data: \n %#v\n", Data())
	fmt.Print("get config example:\n")

	name, ok := GetString("name")
	fmt.Printf("- get string\n ok: %v, val: %v\n", ok, name)

	arr1, ok := GetStringArr("arr1")
	fmt.Printf("- get array\n ok: %v, val: %#v\n", ok, arr1)

	val0, ok := GetString("arr1.0")
	fmt.Printf("- get sub-value by path 'arr.index'\n ok: %v, val: %v\n", ok, val0)

	map1, ok := GetStringMap("map1")
	fmt.Printf("- get map\n ok: %v, val: %#v\n", ok, map1)

	val0, ok = GetString("map1.name")
	fmt.Printf("- get sub-value by path 'map.key'\n ok: %v, val: %v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", DefString("envKey1", ""))

	// Output:
	// get config example:
	// - get string
	// ok: true, val: app2
	// - get array
	// ok: true, val: []string{"alpha", "omega"}
	// - get sub-value by path 'arr.index'
	// ok: true, val: alpha
	// - get map
	// ok: true, val: map[string]string{"name":"inhere", "org":"GitHub"}
	// - get sub-value by path 'map.key'
	// ok: true, val: inhere
	// get env 'envKey' val: /bin/zsh
	// get env 'envKey1' val: defValue
}

func ExampleConfig_DefBool() {
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
	SetEncoder(Yaml, yaml.Encoder)

	buf := new(bytes.Buffer)
	_, err := DumpTo(buf, Yaml)
	if err != nil {
		panic(err)
	}

	fmt.Printf("export config:\n%s", buf.String())

	// Output:
	// arr1:
	// 	- val1
	// 	- val21
	// baseKey: value2
	// debug: false
	// ... ...
}
