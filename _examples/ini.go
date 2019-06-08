// These are some sample code for YAML,TOML,JSON,INI,HCL
package main

import (
	"fmt"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
)

// go run ./examples/ini.go
func main() {
	config.WithOptions(config.ParseEnv)

	// add Decoder and Encoder
	config.AddDriver(ini.Driver)
	// Or
	// config.SetEncoder(config.Ini, ini.Encoder)

	err := config.LoadFiles("testdata/ini_base.ini")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())

	err = config.LoadFiles("testdata/ini_other.ini")
	// config.LoadFiles("testdata/ini_base.ini", "testdata/ini_other.ini")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())
	fmt.Print("get config example:\n")

	name, ok := config.String("name")
	fmt.Printf("- get string\n ok: %v, val: %v\n", ok, name)

	// NOTICE: ini is not support array

	map1, ok := config.StringMap("map1")
	fmt.Printf("- get map\n ok: %v, val: %#v\n", ok, map1)

	val0, ok := config.String("map1.key")
	fmt.Printf("- get sub-value by path 'map.key'\n ok: %v, val: %v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.DefString("envKey1", ""))

	// set value
	config.Set("name", "new name")
	name, ok = config.String("name")
	fmt.Printf("- set string\n ok: %v, val: %v\n", ok, name)

}
