package main

import (
	"github.com/gookit/config"
	"github.com/gookit/config/yaml"
	"fmt"
)

// go run ./examples/yaml.go
func main() {
	config.SetOptions(&config.Options{
		ParseEnv: true,
	})
	// config.SetDecoder(config.Yaml, yaml.Decoder)
	config.SetDriver(config.Yaml, yaml.Decoder, yaml.Encoder)

	err := config.LoadFiles("testdata/yml_base.yml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())

	err = config.LoadFiles("testdata/yml_other.yml")
	// config.LoadFiles("testdata/yml_base.yml", "testdata/yml_other.yml")
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
	fmt.Printf("- get sub-value by path 'arr.index'\n ok: %v, val: %#v\n", ok, val0)

	map1, ok := config.GetStringMap("map1")
	fmt.Printf("- get map\n ok: %v, val: %#v\n", ok, map1)

	val0, ok = config.GetString("map1.key")
	fmt.Printf("- get sub-value by path 'map.key'\n ok: %v, val: %#v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.DefString("envKey1", ""))

	// if you want export config data
	// buf := new(bytes.Buffer)
	// _, err = config.DumpTo(buf, config.Yaml)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("export config:\n%s", buf.String())
}
