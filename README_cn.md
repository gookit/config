# config

[![GoDoc](https://godoc.org/github.com/gookit/config?status.svg)](https://godoc.org/github.com/gookit/config)
[![Build Status](https://travis-ci.org/gookit/config.svg?branch=master)](https://travis-ci.org/gookit/config)
[![Coverage Status](https://coveralls.io/repos/github/gookit/config/badge.svg?branch=master)](https://coveralls.io/github/gookit/config?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/config)](https://goreportcard.com/report/github.com/gookit/config)

golang应用程序配置管理工具库。

> **[EN README](README.md)**

- 支持多种格式: `JSON`(default), `INI`, `YAML`, `TOML`, `HCL`
  - `JSON` 内容支持注释，将自动清除注释
- 支持多个文件/数据加载
- 支持数据覆盖合并，加载多份数据时将按key自动合并
- 支持从远程URL加载配置数据
- 支持从命令行参数设置配置数据
- 支持按路径获取子级值。 e.g `map.key` `arr.2`
- 支持解析ENV变量名称。 like `shell: ${SHELL}` -> `shell: /bin/zsh`
- 简洁的使用API `Get` `Int` `String` `Bool` `Ints` `IntMap` `Strings` `StringMap` ...
- 完善的单元测试(coverage > 90%)

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
	// 设置选项支持 ENV 解析
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

### 获取数据

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

## API方法参考

### 载入配置

- `LoadData(dataSource ...interface{}) (err error)`
- `LoadFlags(keys []string) (err error)`
- `LoadExists(sourceFiles ...string) (err error)`
- `LoadFiles(sourceFiles ...string) (err error)`
- `LoadRemote(format, url string) (err error)`
- `LoadSources(format string, src []byte, more ...[]byte) (err error)`
- `LoadStrings(format string, str string, more ...string) (err error)`

### 获取值

- `Bool(key string) (value bool, ok bool)`
- `Int(key string) (value int, ok bool)`
- `Int64(key string) (value int64, ok bool)`
- `Ints(key string) (arr []int, ok bool)`
- `IntMap(key string) (mp map[string]int, ok bool)`
- `Float(key string) (value float64, ok bool)`
- `String(key string) (value string, ok bool)`
- `Strings(key string) (arr []string, ok bool)`
- `StringMap(key string) (mp map[string]string, ok bool)`
- `Get(key string, findByPath ...bool) (value interface{}, ok bool)`

### 设置值

- `Set(key string, val interface{}, setByPath ...bool) (err error)`

## 单元测试

```bash
go test -cover
// contains all sub-folder
go test -cover ./...
```

## 相关包

- Ini 解析 [gookit/ini/parser](https://github.com/gookit/ini/tree/master/parser)
- Yaml 解析 [go-yaml](https://github.com/go-yaml/yaml)
- Toml 解析 [go toml](https://github.com/BurntSushi/toml)
- 数据合并 [mergo](https://github.com/imdario/mergo)

## License

**MIT**
