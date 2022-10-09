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

	name := config.String("name")
	fmt.Printf("- get string\n val: %v\n", name)

	arr1 := config.Strings("arr1")
	fmt.Printf("- get array\n val: %#v\n", arr1)

	val0 := config.String("arr1.0")
	fmt.Printf("- get sub-value by path 'arr.index'\n val: %v\n", val0)

	map1 := config.StringMap("map1")
	fmt.Printf("- get map\n val: %#v\n", map1)

	val0 = config.String("map1.name")
	fmt.Printf("- get sub-value by path 'map.key'\n val: %v\n", val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.String("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.String("envKey1", ""))

}
