# config

[![GoDoc](https://godoc.org/github.com/gookit/config?status.svg)](https://godoc.org/github.com/gookit/config)

golang应用程序配置管理工具库。

- 支持多种格式: `JSON`(default), `INI`, `YAML`, `TOML`, `HCL`
- 支持多文件/数据加载
- 支持数据覆盖合并
- 支持按路径获取子级值, e.g `map.key` `arr.2`
- 支持解析ENV变量名称. like `envKey: ${SHELL}` -> `envKey: /bin/zsh`
- 通用的使用API `Get` `Int` `String` `Bool` `Ints` `IntMap` `Strings` `StringMap` ...

> **[EN README](README.md)**

## Godoc

- [godoc for gopkg](https://godoc.org/gopkg.in/gookit/config.v1)
- [godoc for github](https://godoc.org/github.com/gookit/config)

## 快速使用

这里使用yaml格式作为示例(`testdata/yml_other.yml`):

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

- 使用, demo请看 [examples/yaml.go](examples/yaml.go):

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
		ParseEnv: true, // 启用解析ENV变量
	})
	
	// 除了json格式是内置的外，其他格式都需在使用前添加解析驱动
	// config.SetDecoder(config.Yaml, yaml.Decoder)
	config.AddDriver(config.Yaml, yaml.Driver)
	// OR
	// config.DecoderEncoder(config.Yaml, yaml.Decoder, yaml.Encoder)

    // 加载配置，可以同时传入多个文件
	err := config.LoadFiles("testdata/yml_base.yml")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", config.Data())

	// 加载更多文件
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

- 输出:

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

## 有用的包

### Yaml 解析

- [go-yaml](https://github.com/go-yaml/yaml) yaml parser

### Toml 解析

- [go toml](https://github.com/BurntSushi/toml) toml parser

### Ini 解析

- [gookit/ini/parser](https://github.com/gookit/ini/parser) ini parser

### 数据合并

- [mergo](https://github.com/imdario/mergo) merge data

## 其他

使用INI作为简单的配置管理

- [gookit/ini](https://github.com/gookit/ini) 

## License

MIT
