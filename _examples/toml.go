package main

import (
	"fmt"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/toml"
)

// go run ./examples/toml.go
func main() {
	config.WithOptions(config.ParseEnv)

	// add Decoder and Encoder
	config.AddDriver(toml.Driver)

	err := config.LoadFiles("testdata/toml_base.toml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())

	err = config.LoadFiles("testdata/toml_other.toml")
	// config.LoadFiles("testdata/toml_base.toml", "testdata/toml_other.toml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())
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

}
