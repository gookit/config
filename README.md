# config

[![GoDoc](https://godoc.org/github.com/gookit/config?status.svg)](https://godoc.org/github.com/gookit/config)

golang application config manage implement. 

- generic api `Get` `GetInt` `GetString` `GetBool` `GetStringArr` ...
- support file format: `json`(default), `yaml`, `toml`
- support multi file/data load
- support data override merge
- support get sub value by path, like `map.key` `arr.2`

## Godoc

- [godoc for gopkg](https://godoc.org/gopkg.in/gookit/config.v1)
- [godoc for github](https://godoc.org/github.com/gookit/config)

## Usage

Here using the yaml format as an example(`testdata/yml_other.yml`):

```yaml
name: app2
debug: false
baseKey: value2

map1:
    key: val2
    key2: val20

arr1:
    - val1
    - val21
```

- usage:

```go
package main

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
```

- output:

```text
get 'name', ok: true, val: "app2"
get 'arr1', ok: true, val: []string{"val1", "val21"}
get sub 'arr1.0', ok: true, val: "val1"
get 'map1', ok: true, val: map[string]string{"key":"val2", "key2":"val20"}
get sub 'map1.key', ok: true, val: "val2"
```

## Useful packages

### ini config use

- [go-ini/ini](https://github.com/go-ini/ini) ini parser and config manage
- [dombenson/go-ini](https://github.com/dombenson/go-ini) ini parser and config manage

### yaml

- [go-yaml](https://github.com/go-yaml/yaml) ini parser


### toml

- [go toml](https://github.com/BurntSushi/toml) toml parser

### data merge

- [mergo](github.com/imdario/mergo) merge data

## License

MIT
