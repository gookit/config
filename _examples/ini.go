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

	name := config.String("name")
	fmt.Printf("- get string\n val: %v\n", name)

	// NOTICE: ini is not support array

	map1 := config.StringMap("map1")
	fmt.Printf("- get map\n val: %#v\n", map1)

	val0 := config.String("map1.key")
	fmt.Printf("- get sub-value by path 'map.key'\n val: %v\n", val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.String("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.String("envKey1", ""))

	// set value
	_ = config.Set("name", "new name")
	name = config.String("name")
	fmt.Printf("- set string\n val: %v\n", name)

}
