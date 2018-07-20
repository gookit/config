package toml

import (
	"fmt"
	"github.com/gookit/config"
)

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

func Example() {
	config.SetOptions(&config.Options{
		ParseEnv: true,
	})

	// add Decoder and Encoder
	config.AddDriver(config.Toml, Driver)
	// Or
	// config.DecoderEncoder(config.Toml, Decoder, Encoder)

	err := config.LoadFiles("testdata/toml_base.toml")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", Data())

	// load more files
	err = config.LoadFiles("testdata/toml_other.toml")
	// can also
	// config.LoadFiles("testdata/toml_base.toml", "testdata/toml_other.toml")
	if err != nil {
		panic(err)
	}

	// load from string
	config.LoadSources(config.Toml, []byte(tomlStr))

	// fmt.Printf("config data: \n %#v\n", Data())
	fmt.Print("get config example:\n")

	name, ok := config.String("name")
	fmt.Printf("- get string\n ok: %v, val: %v\n", ok, name)

	arr1, ok := config.Strings("arr1")
	fmt.Printf("- get array\n ok: %v, val: %#v\n", ok, arr1)

	val0, ok := config.String("arr1.0")
	fmt.Printf("- get sub-value by path 'arr.index'\n ok: %v, val: %v\n", ok, val0)

	map1, ok := config.StringMap("map1")
	fmt.Printf("- get map\n ok: %v, val: %#v\n", ok, map1)

	val0, ok = config.String("map1.name")
	fmt.Printf("- get sub-value by path 'map.key'\n ok: %v, val: %v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.DefString("envKey1", ""))

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
