package yaml

import (
	"fmt"
	"github.com/gookit/config"
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
	// only add decoder
	// config.SetDecoder(config.Yaml, Decoder)
	// Or
	config.AddDriver(config.Yaml, Driver)
	// Or
	// config.DecoderEncoder(config.Yaml, Decoder, Encoder)

	err := config.LoadFiles("testdata/yml_other.yml")
	if err != nil {
		panic(err)
	}

	// load from string
	config.LoadSources(config.Yaml, []byte(yamlStr))

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

	// Output:
	// get config example:
	// - get string
	// ok: true, val: app2
	// - get array
	// ok: true, val: []string{"val1", "val21"}
	// - get sub-value by path 'arr.index'
	// ok: true, val: "val1"
	// - get map
	// ok: true, val: map[string]string{"key":"val2", "key2":"val20"}
	// - get sub-value by path 'map.key'
	// ok: true, val: "val2"
	// get env 'envKey' val: /bin/zsh
	// get env 'envKey1' val: defValue
}
