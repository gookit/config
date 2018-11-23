# config

[![GoDoc](https://godoc.org/github.com/gookit/config?status.svg)](https://godoc.org/github.com/gookit/config)
[![Build Status](https://travis-ci.org/gookit/config.svg?branch=master)](https://travis-ci.org/gookit/config)
[![Coverage Status](https://coveralls.io/repos/github/gookit/config/badge.svg?branch=master)](https://coveralls.io/github/gookit/config?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/config)](https://goreportcard.com/report/github.com/gookit/config)

golang application config manage tool library. 

- support multi format: `JSON`(default), `INI`, `YAML`, `TOML`, `HCL`
  - `JSON` content support comments. will auto clear comments
- support multi file/data load
- support for loading configuration data from remote URLs
- support for setting configuration data from command line arguments
- support data overlay and merge, automatically load by key when loading multiple copies of data
- support get sub value by path, like `map.key` `arr.2`
- support parse ENV name. like `envKey: ${SHELL}` -> `envKey: /bin/zsh`
- generic api `Get` `Int` `String` `Bool` `Ints` `IntMap` `Strings` `StringMap` ...
- complete unit test(coverage > 90%)

> **[中文说明](README_cn.md)**

## Only use INI

> If you just want to use INI for simple config management, recommended use [gookit/ini](https://github.com/gookit/ini)

## Godoc

- [godoc for gopkg](https://godoc.org/gopkg.in/gookit/config.v1)
- [godoc for github](https://godoc.org/github.com/gookit/config)

## Usage

Here using the yaml format as an example(`testdata/yml_other.yml`):

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

### Load data

> examples code please see [_examples/yaml.go](_examples/yaml.go):

```go
package main

import (
    "github.com/gookit/config"
    "github.com/gookit/config/yaml"
)

// go run ./examples/yaml.go
func main() {
	config.WithOptions(config.ParseEnv)
	
	// add driver for support yaml content
	config.AddDriver(yaml.Driver)
	// config.SetDecoder(config.Yaml, yaml.Decoder)

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
}
```

### Read data

- get integer

```go
age, ok := config.Int("age")
fmt.Print(ok, age) // true 100
```

- get bool

```go
val, ok := config.Bool("debug")
fmt.Print(ok, age) // true true
```

- get string

```go
name, ok := config.String("name")
fmt.Print(ok, name) // true inhere
```

- get strings(slice)

```go
arr1, ok := config.Strings("arr1")
fmt.Printf("%v %#v", ok, arr1) // true []string{"val1", "val21"}
```

- get string map

```go
val, ok := config.StringMap("map1")
fmt.Printf("%v %#v",ok, val) // true map[string]string{"key":"val2", "key2":"val20"}
```

- value contains ENV var

```go
value, ok := config.String("shell")
fmt.Print(ok, value) // true /bin/zsh
```

- get value by key path

```go
// from array
value, ok := config.String("arr1.0")
fmt.Print(ok, value) // true "val1"

// from map
value, ok := config.String("map1.key")
fmt.Print(ok, value) // true "val2"
```

- setting new value

```go
// set value
config.Set("name", "new name")
name, ok = config.String("name")
fmt.Print(ok, name) // true "new name"
```

## API Methods Refer

### Load Config

- `LoadData(dataSource ...interface{}) (err error)`
- `LoadFlags(keys []string) (err error)`
- `LoadExists(sourceFiles ...string) (err error)`
- `LoadFiles(sourceFiles ...string) (err error)`
- `LoadRemote(format, url string) (err error)`
- `LoadSources(format string, src []byte, more ...[]byte) (err error)`
- `LoadStrings(format string, str string, more ...string) (err error)`

### Getting Values

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

### Setting Values

- `Set(key string, val interface{}, setByPath ...bool) (err error)`

## Run Tests

```bash
go test -cover
// contains all sub-folder
go test -cover ./...
```

## Related Packages

- Ini parse [gookit/ini/parser](https://github.com/gookit/ini/tree/master/parser)
- Yaml parse [go-yaml](https://github.com/go-yaml/yaml)
- Toml parse [go toml](https://github.com/BurntSushi/toml)
- Data merge [mergo](https://github.com/imdario/mergo)

### Ini Config Use

- [gookit/ini](https://github.com/gookit/ini) ini config manage

## License

**MIT**
