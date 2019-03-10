# config

[![GoDoc](https://godoc.org/github.com/gookit/config?status.svg)](https://godoc.org/github.com/gookit/config)
[![Build Status](https://travis-ci.org/gookit/config.svg?branch=master)](https://travis-ci.org/gookit/config)
[![Coverage Status](https://coveralls.io/repos/github/gookit/config/badge.svg?branch=master)](https://coveralls.io/github/gookit/config?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/config)](https://goreportcard.com/report/github.com/gookit/config)

Golang application config manage tool library. 

> **[中文说明](README_cn.md)**

## Features

- Support multi format: `JSON`(default), `INI`, `YAML`, `TOML`, `HCL`, `ENV`, `Flags`
  - `JSON` content support comments. will auto clear comments
- Support multi-file and multi-data loading
- Support loading configuration from os ENV
- Support for loading configuration data from remote URLs
- Support for setting configuration data from command line arguments(`flags`)
- Support data overlay and merge, automatically load by key when loading multiple copies of data
- Support get sub value by path, like `map.key` `arr.2`
- Support parse ENV name. like `envKey: ${SHELL}` -> `envKey: /bin/zsh`
- Generic api `Get` `Int` `Uint` `Int64` `Float` `String` `Bool` `Ints` `IntMap` `Strings` `StringMap` ...
- Complete unit test(code coverage > 95%)

> Provide a sub-package `dotenv` that supports importing data from files (eg `.env`) to ENV

## Only use INI

> If you just want to use INI for simple config management, recommended use [gookit/ini](https://github.com/gookit/ini)

## GoDoc

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

	// load more files
	err = config.LoadFiles("testdata/yml_other.yml")
	// can also load multi at once
	// err := config.LoadFiles("testdata/yml_base.yml", "testdata/yml_other.yml")
	if err != nil {
		panic(err)
	}
	
	// fmt.Printf("config data: \n %#v\n", config.Data())
}
```

### Read data

- get integer

```go
age := config.Int("age")
fmt.Print(age) // 100
```

- Get bool

```go
val := config.Bool("debug")
fmt.Print(val) // true
```

- Get string

```go
name := config.String("name")
fmt.Print(name) // inhere
```

- Get strings(slice)

```go
arr1 := config.Strings("arr1")
fmt.Printf("%#v", arr1) // []string{"val1", "val21"}
```

- Get string map

```go
val := config.StringMap("map1")
fmt.Printf("%#v",val) // map[string]string{"key":"val2", "key2":"val20"}
```

- Value contains ENV var

```go
value := config.String("shell")
fmt.Print(value) // "/bin/zsh"
```

- Get value by key path

```go
// from array
value := config.String("arr1.0")
fmt.Print(value) // "val1"

// from map
value := config.String("map1.key")
fmt.Print(value) // "val2"
```

- Setting new value

```go
// set value
config.Set("name", "new name")
name = config.String("name")
fmt.Print(name) // "new name"
```

## Load from flags

> Support simple flags parameter parsing, loading

```go
// flags like: --name inhere --env dev --age 99 --debug

// load flag info
keys := []string{"name", "env", "age:int" "debug:bool"}
err := config.LoadFlags(keys)

// read
config.String("name") // "inhere"
config.String("env") // "dev"
config.Int("age") // 99
config.Bool("debug") // true
```

## Load from ENV

```go
// os env: APP_NAME=config APP_DEBUG=true
// load ENV info
config.LoadOSEnv([]string{"app_name", "app_debug"})

// read
config.Bool("app_debug") // true
config.String("app_name") // "config"
```

## API Methods Refer

### Load Config

- `LoadOSEnv(keys []string)` Load from os ENV
- `LoadData(dataSource ...interface{}) (err error)` Load from struts or maps
- `LoadFlags(keys []string) (err error)` Load from CLI flags
- `LoadExists(sourceFiles ...string) (err error)` 
- `LoadFiles(sourceFiles ...string) (err error)`
- `LoadRemote(format, url string) (err error)`
- `LoadSources(format string, src []byte, more ...[]byte) (err error)`
- `LoadStrings(format string, str string, more ...string) (err error)`

### Getting Values

- `Bool(key string, defVal ...bool) bool`
- `Int(key string, defVal ...int) int`
- `Uint(key string, defVal ...uint) uint`
- `Int64(key string, defVal ...int64) int64`
- `Ints(key string) (arr []int)`
- `IntMap(key string) (mp map[string]int)`
- `Float(key string, defVal ...float64) float64`
- `String(key string, defVal ...string) string`
- `Strings(key string) (arr []string)`
- `StringMap(key string) (mp map[string]string)`
- `Get(key string, findByPath ...bool) (value interface{})`

### Setting Values

- `Set(key string, val interface{}, setByPath ...bool) (err error)`

### Useful Methods

- `Getenv(name string, defVal ...string) (val string)`
- `AddDriver(driver Driver)`
- `Data() map[string]interface{}`
- `Exists(key string, findByPath ...bool) bool`
- `DumpTo(out io.Writer, format string) (n int64, err error)`

## Run Tests

```bash
go test -cover
// contains all sub-folder
go test -cover ./...
```

## Gookit packages

- [gookit/ini](https://github.com/gookit/ini) Go config management, use INI files
- [gookit/rux](https://github.com/gookit/rux) Simple and fast request router for golang HTTP 
- [gookit/gcli](https://github.com/gookit/gcli) build CLI application, tool library, running CLI commands
- [gookit/event](https://github.com/gookit/event) Lightweight event manager and dispatcher implements by Go
- [gookit/cache](https://github.com/gookit/cache) Generic cache use and cache manager for golang. support File, Memory, Redis, Memcached.
- [gookit/config](https://github.com/gookit/config) Go config management. support JSON, YAML, TOML, INI, HCL, ENV and Flags
- [gookit/color](https://github.com/gookit/color) A command-line color library with true color support, universal API methods and Windows support
- [gookit/filter](https://github.com/gookit/filter) Provide filtering, sanitizing, and conversion of golang data
- [gookit/validate](https://github.com/gookit/validate) Use for data validation and filtering. support Map, Struct, Form data
- [gookit/goutil](https://github.com/gookit/goutil) Some utils for the Go: string, array/slice, map, format, cli, env, filesystem, test and more
- More please see https://github.com/gookit

## See also

- Ini parse [gookit/ini/parser](https://github.com/gookit/ini/tree/master/parser)
- Ini config [gookit/ini](https://github.com/gookit/ini)
- Yaml parse [go-yaml](https://github.com/go-yaml/yaml)
- Toml parse [go toml](https://github.com/BurntSushi/toml)
- Data merge [mergo](https://github.com/imdario/mergo)

## License

**MIT**
