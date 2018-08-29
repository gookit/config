# config

[![GoDoc](https://godoc.org/github.com/gookit/config?status.svg)](https://godoc.org/github.com/gookit/config)
[![Build Status](https://travis-ci.org/gookit/config.svg?branch=master)](https://travis-ci.org/gookit/config)
[![Coverage Status](https://coveralls.io/repos/github/gookit/config/badge.svg?branch=master)](https://coveralls.io/github/gookit/config?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/config)](https://goreportcard.com/report/github.com/gookit/config)

golang应用程序配置管理工具库。

- 支持多种格式: `JSON`(default), `INI`, `YAML`, `TOML`, `HCL`
- 支持多个文件/数据加载
- 支持数据覆盖合并，将按key自动合并
- 支持按路径获取子级值。 e.g `map.key` `arr.2`
- 支持解析ENV变量名称。 like `shell: ${SHELL}` -> `shell: /bin/zsh`
- 简洁的使用API `Get` `Int` `String` `Bool` `Ints` `IntMap` `Strings` `StringMap` ...
- 完善的单元测试(coverage > 90%)

> **[EN README](README.md)**

## 只使用INI

> 如果你仅仅想用INI来做简单配置管理，推荐使用 [gookit/ini](https://github.com/gookit/ini)

## Godoc

- [godoc for gopkg](https://godoc.org/gopkg.in/gookit/config.v1)
- [godoc for github](https://godoc.org/github.com/gookit/config)

## 快速使用

这里使用yaml格式作为示例(`testdata/yml_other.yml`):

```yaml
name: app2
debug: false
baseKey: value2
shell: ${SHELL}
envKey1: ${NotExist|defValue}

map1:
    key: val2
    key2: val20

arr1:
    - val1
    - val21
```

### 载入数据

> 示例代码请看 [_examples/yaml.go](_examples/yaml.go):

```go
package main

import (
    "github.com/gookit/config"
    "github.com/gookit/config/yaml"
)

// go run ./examples/yaml.go
func main() {
	config.WithOptions(config.ParseEnv)

    // 添加驱动程序以支持yaml内容解析（除了JSON是默认支持，其他的则是按需使用）
    config.AddDriver(yaml.Driver)

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
}
```

### 读取数据

- 获取整型

```go
age, ok := config.Int("age")
fmt.Print(ok, age) // true 100
```

- 获取布尔值

```go
val, ok := config.Bool("debug")
fmt.Print(ok, age) // true true
```

- 获取字符串

```go
name, ok := config.String("name")
fmt.Print(ok, name) // true inhere
```

- 获取字符串数组

```go
arr1, ok := config.Strings("arr1")
fmt.Printf("%v %#v", ok, arr1) // true []string{"val1", "val21"}
```

- 获取字符串KV映射

```go
val, ok := config.StringMap("map1")
fmt.Printf("%v %#v",ok, val) // true map[string]string{"key":"val2", "key2":"val20"}
```

- 值包含ENV变量

```go
value, ok := config.String("shell")
fmt.Print(ok, value) // true /bin/zsh
```

- 通过key路径获取值

```go
// from array
value, ok := config.String("arr1.0")
fmt.Print(ok, value) // true "val1"

// from map
value, ok := config.String("map1.key")
fmt.Print(ok, value) // true "val2"
```

- 设置新的值

```go
// set value
config.Set("name", "new name")
name, ok = config.String("name")
fmt.Print(ok, name) // true "new name"
```

## 单元测试

```bash
go test -cover
// contains all sub-folder
go test -cover ./...
```

## 有用的包

- Ini 解析 [gookit/ini/parser](https://github.com/gookit/ini/tree/master/parser)
- Yaml 解析 [go-yaml](https://github.com/go-yaml/yaml)
- Toml 解析 [go toml](https://github.com/BurntSushi/toml)
- 数据合并 [mergo](https://github.com/imdario/mergo)

## 其他

使用INI作为简单的配置管理

- [gookit/ini](https://github.com/gookit/ini) 

## License

**MIT**
