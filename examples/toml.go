package main

import (
	"fmt"
	"github.com/gookit/config"
	"github.com/gookit/config/toml"
)

// go run ./examples/toml.go
func main()  {
	config.SetOptions(&config.Options{
		ParseEnv: true,
	})
	config.SetDriver(config.Toml, toml.Decoder, toml.Encoder)

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

	name, ok := config.GetString("name")
	fmt.Printf("- get string\n ok: %v, val: %v\n", ok, name)

	arr1, ok := config.GetStringArr("arr1")
	fmt.Printf("- get array\n ok: %v, val: %#v\n", ok, arr1)

	val0, ok := config.GetString("arr1.0")
	fmt.Printf("- get sub-value by path 'arr.index'\n ok: %v, val: %v\n", ok, val0)

	map1, ok := config.GetStringMap("map1")
	fmt.Printf("- get map\n ok: %v, val: %#v\n", ok, map1)

	val0, ok = config.GetString("map1.name")
	fmt.Printf("- get sub-value by path 'map.key'\n ok: %v, val: %v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.DefString("envKey1", ""))

}
