package main

import (
	"fmt"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
)

// go run ./examples/yaml.go
func main() {
	config.WithOptions(config.ParseEnv)

	// only add decoder
	// config.SetDecoder(config.Yaml, yamlv3.Decoder)
	// Or
	config.AddDriver(yamlv3.Driver)

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

	name := config.String("name")
	fmt.Printf("- get string\n val: %v\n", name)

	arr1 := config.Strings("arr1")
	fmt.Printf("- get array\n val: %#v\n", arr1)

	val0 := config.String("arr1.0")
	fmt.Printf("- get sub-value by path 'arr.index'\n val: %#v\n", val0)

	map1 := config.StringMap("map1")
	fmt.Printf("- get map\n val: %#v\n", map1)

	val0 = config.String("map1.key")
	fmt.Printf("- get sub-value by path 'map.key'\n val: %#v\n", val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.String("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.String("envKey1", ""))

	// set value
	config.Set("name", "new name")
	name = config.String("name")
	fmt.Printf("- set string\n  val: %v\n", name)

	// if you want export config data
	// buf := new(bytes.Buffer)
	// _, err = config.DumpTo(buf, config.Yaml)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("export config:\n%s", buf.String())
}
