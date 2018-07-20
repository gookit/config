# config

[![GoDoc](https://godoc.org/github.com/gookit/config?status.svg)](https://godoc.org/github.com/gookit/config)
[![cover.run](https://cover.run/go/https:/github.com/gookit/config.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=https%3A%2Fgithub.com%2Fgookit%2Fconfig)

golang application config manage tool library. 

- support multi format: `JSON`(default), `INI`, `YAML`, `TOML`, `HCL`
- support multi file/data load
- support data override merge
- support get sub value by path, like `map.key` `arr.2`
- support parse env name. like `envKey: ${SHELL}` -> `envKey: /bin/zsh`
- generic api `Get` `Int` `String` `Bool` `Ints` `IntMap` `Strings` `StringMap` ...

> **[中文说明](README_cn.md)**

## Godoc

- [godoc for gopkg](https://godoc.org/gopkg.in/gookit/config.v1)
- [godoc for github](https://godoc.org/github.com/gookit/config)

## Usage

Here using the yaml format as an example(`testdata/yml_other.yml`):

```yaml
name: app2
debug: false
baseKey: value2
envKey: ${SHELL}
envKey1: ${NotExist|defValue}

map1:
    key: val2
    key2: val20

arr1:
    - val1
    - val21
```

- usage, see [examples/yaml.go](examples/yaml.go):

```go
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
	config.DecoderEncoder(config.Yaml, yaml.Decoder, yaml.Encoder)

	err := config.LoadFiles("testdata/yml_base.yml")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", config.Data())

	// load more files
	err = config.LoadFiles("testdata/yml_other.yml")
	// can also load multi at once
	// err := config.LoadFiles("testdata/yml_base.yml", "testdata/yml_other.yml")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", config.Data())
	fmt.Print("get config example:\n")

	name, ok := config.String("name")
	fmt.Printf("get string\n - ok: %v, val: %v\n", ok, name)

	arr1, ok := config.Strings("arr1")
	fmt.Printf("get array\n - ok: %v, val: %#v\n", ok, arr1)

	val0, ok := config.String("arr1.0")
	fmt.Printf("get sub-value by path 'arr.index'\n - ok: %v, val: %#v\n", ok, val0)

	map1, ok := config.StringMap("map1")
	fmt.Printf("get map\n - ok: %v, val: %#v\n", ok, map1)

	val0, ok = config.String("map1.key")
	fmt.Printf("get sub-value by path 'map.key'\n - ok: %v, val: %#v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.DefString("envKey1", ""))

	// set value
	config.Set("name", "new name")
	name, ok = config.String("name")
	fmt.Printf("set string\n - ok: %v, val: %v\n", ok, name)
	
	// if you want export config data
	// buf := new(bytes.Buffer)
	// _, err = config.DumpTo(buf, config.Yaml)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("export config:\n%s", buf.String())
}
```

- output:

```text
get config example:
get string
 - ok: true, val: app2
get array
 - ok: true, val: []string{"val1", "val21"}
get sub-value by path 'arr.index'
 - ok: true, val: "val1"
get map
 - ok: true, val: map[string]string{"key":"val2", "key2":"val20"}
get sub-value by path 'map.key'
 - ok: true, val: "val2"
get env 'envKey' val: /bin/zsh
get env 'envKey1' val: defValue
set string
- ok: true, val: new name
```

## Useful packages

### Yaml parse

- [go-yaml](https://github.com/go-yaml/yaml) yaml parser

### Toml parse

- [go toml](https://github.com/BurntSushi/toml) toml parser

### Ini parse

- [gookit/ini/parser](https://github.com/gookit/ini/parser) ini parser

### Data merge

- [mergo](https://github.com/imdario/mergo) merge data

### Ini config use

- [gookit/ini](https://github.com/gookit/ini/parser) ini config manage

## License

MIT
