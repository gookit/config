package main

import (
	"github.com/gookit/config"
	"github.com/gookit/config/yaml"
	"fmt"
)

func main() {
	config.SetOptions(&config.Options{})
	config.SetDecoder(config.Yaml, yaml.Decoder)

	config.LoadFiles("testdata/yml_base.yml")

	fmt.Printf("config data: \n %#v\n", config.Data())

	config.LoadFiles("testdata/yml_other.yml")
	// config.LoadFiles("testdata/yml_base.yml", "testdata/yml_other.yml")

	fmt.Printf("config data: \n %#v\n", config.Data())

	name, ok := config.GetString("name")
	fmt.Printf("get 'name', ok: %v, val: %#v\n", ok, name)

	arr1, ok := config.GetStringArr("arr1")
	fmt.Printf("get 'arr1', ok: %v, val: %#v\n", ok, arr1)

	val0, ok := config.GetString("arr1.0")
	fmt.Printf("get sub 'arr1.0', ok: %v, val: %#v\n", ok, val0)

	map1, ok := config.GetStringMap("map1")
	fmt.Printf("get 'map1', ok: %v, val: %#v\n", ok, map1)

	val0, ok = config.GetString("map1.key")
	fmt.Printf("get sub 'map1.key', ok: %v, val: %#v\n", ok, val0)
}
