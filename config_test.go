package config

import (
	"fmt"
	"github.com/gookit/config/yaml"
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

func Example() {
	// add yaml decoder
	SetDecoder(Yaml, yaml.Decoder)
	LoadFiles("testdata/yml_other.yml")

	name, ok := GetString("name")
	fmt.Printf("get 'name', ok: %v, val: %#v\n", ok, name)

	arr1, ok := GetStringArr("arr1")
	fmt.Printf("get 'arr1', ok: %v, val: %#v\n", ok, arr1)

	val0, ok := GetString("arr1.0")
	fmt.Printf("get sub 'arr1.0', ok: %v, val: %#v\n", ok, val0)

	map1, ok := GetStringMap("map1")
	fmt.Printf("get 'map1', ok: %v, val: %#v\n", ok, map1)

	val0, ok = GetString("map1.key")
	fmt.Printf("get sub 'map1.key', ok: %v, val: %#v\n", ok, val0)

	// Output:
	// get 'name', ok: true, val: "app2"
	// get 'arr1', ok: true, val: []string{"val1", "val21"}
	// get sub 'arr1.0', ok: true, val: "val1"
	// get 'map1', ok: true, val: map[string]string{"key":"val2", "key2":"val20"}
	// get sub 'map1.key', ok: true, val: "val2"
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
