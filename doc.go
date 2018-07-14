/*
golang application config manage implement.

Source code and other details for the project are available at GitHub:

   https://github.com/gookit/config

Here using the yaml format as an example(yml_other.yml):

	name: app2
	debug: false
	baseKey: value2

	map1:
	    key: val2
	    key2: val20

	arr1:
	    - val1
	    - val21

usage:

	import (
		"github.com/gookit/config"
		"github.com/gookit/config/yaml"
		"fmt"
	)

	// add yaml decoder
	config.SetDecoder(config.Yaml, yaml.Decoder)
	config.LoadFiles("testdata/yml_other.yml")

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

output:

	get 'name', ok: true, val: "app2"
	get 'arr1', ok: true, val: []string{"val1", "val21"}
	get sub 'arr1.0', ok: true, val: "val1"
	get 'map1', ok: true, val: map[string]string{"key":"val2", "key2":"val20"}
	get sub 'map1.key', ok: true, val: "val2"

 */
package config
