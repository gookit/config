# Config

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/config?style=flat-square)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/d6ac163ee63649ec92c1566e42f09c11)](https://app.codacy.com/app/inhere/config)
[![GoDoc](https://godoc.org/github.com/gookit/config?status.svg)](https://godoc.org/github.com/gookit/config)
[![Build Status](https://travis-ci.org/gookit/config.svg?branch=master)](https://travis-ci.org/gookit/config)
[![Actions Status](https://github.com/gookit/config/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/config/actions)
[![Coverage Status](https://coveralls.io/repos/github/gookit/config/badge.svg?branch=master)](https://coveralls.io/github/gookit/config?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/config)](https://goreportcard.com/report/github.com/gookit/config)

简洁、功能完善的Golang应用程序配置管理工具库

> **[EN README](README.md)**

## 功能简介

- 支持多种格式: `JSON`(默认), `INI`, `YAML`, `TOML`, `HCL`, `ENV`, `Flags`
  - `JSON` 内容支持注释，可以设置解析时清除注释
  - 其他驱动都是按需使用，不使用的不会加载编译到应用中
- 支持多个文件、多数据加载
- 支持从 OS ENV 变量数据加载配置
- 支持从远程 URL 加载配置数据
- 支持从命令行参数(`flags`)设置配置数据
- 支持数据覆盖合并，加载多份数据时将按key自动合并
- 支持将全部或部分配置数据绑定到结构体 `config.BindStruct("key", &s)`
- 支持通过 `.` 分隔符来按路径获取子级值，也支持自定义分隔符。 e.g `map.key` `arr.2`
- 支持解析ENV变量名称。 like `shell: ${SHELL}` -> `shell: /bin/zsh`
- 简洁的使用API `Get` `Int` `Uint` `Int64` `String` `Bool` `Ints` `IntMap` `Strings` `StringMap` ...
- 完善的单元测试(code coverage > 95%)

> 提供一个子包 `dotenv`，支持从文件（eg `.env`）中导入数据到ENV

## 只使用INI

> 如果你仅仅想用INI来做简单配置管理，推荐使用 [gookit/ini](https://github.com/gookit/ini)

## GoDoc

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
    "github.com/gookit/config/v2"
    "github.com/gookit/config/v2/yaml"
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

### 绑定数据到结构体

> 注意：结构体默认的绑定映射tag是 `mapstructure`，可以通过设置 `Options.TagName` 来更改它

```go
user := struct {
    Age  int
    Kye  string
    UserName string `mapstructure:"user_name"`
    Tags []int
}{}
err = config.BindStruct("user", &user)

fmt.Println(user.UserName) // inhere
```

**更改结构标签名称**

```go
config.WithOptions(func(opt *Options) {
    opt.TagName = "config"
})
```

可以使用空字符串将所有配置数据绑定到结构:

```go
config.BindStruct("", &myConf)
```

### 快速获取数据

```go
// 获取整型
age := config.Int("age")
fmt.Print(age) // 100

// 获取布尔值
val := config.Bool("debug")
fmt.Print(val) // true

// 获取字符串
name := config.String("name")
fmt.Print(name) // inhere

// 获取字符串数组
arr1 := config.Strings("arr1")
fmt.Printf("%v %#v", arr1) // []string{"val1", "val21"}

// 获取字符串KV映射
val := config.StringMap("map1")
fmt.Printf("%v %#v",val) // map[string]string{"key":"val2", "key2":"val20"}

// 值包含ENV变量
value := config.String("shell")
fmt.Print(value) // /bin/zsh

// 通过key路径获取值
// from array
value := config.String("arr1.0")
fmt.Print(value) // "val1"

// from map
value := config.String("map1.key")
fmt.Print(value) // "val2"
```

### 设置新的值

```go
// set value
config.Set("name", "new name")
// get
name = config.String("name")
fmt.Print(name) // new name
```

## 从ENV载入数据

```go
// os env: APP_NAME=config APP_DEBUG=true
// load ENV info
config.LoadOSEnv([]string{"APP_NAME", "APP_NAME"}, true)

// read
config.Bool("app_debug") // true
config.String("app_name") // "config"
```

## 从命令行参数载入数据

> 支持简单的命令行 `flag` 参数解析，加载数据

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

## 创建自定义实例

您可以创建自定义配置实例：

```go
// create new instance, will auto register JSON driver
myConf := config.New("my-conf")

// create empty instance
myConf := config.NewEmpty("my-conf")

// create and with some options
myConf := config.NewWithOptions("my-conf", config.ParseEnv, config.ReadOnly)
```

## 可用选项

```go
// Options config options
type Options struct {
	// parse env value. like: "${EnvName}" "${EnvName|default}"
	ParseEnv bool
	// config is readonly. default is False
	Readonly bool
	// enable config data cache. default is False
	EnableCache bool
	// parse key, allow find value by key path. default is True eg: 'key.sub' will find `map[key]sub`
	ParseKey bool
	// tag name for binding data to struct
	TagName string
	// the delimiter char for split key, when `FindByPath=true`. default is '.'
	Delimiter byte
	// default write format. default is JSON
	DumpFormat string
	// default input format. default is JSON
	ReadFormat string
}
```

## API方法参考

### 载入配置

- `LoadOSEnv(keys []string)` 从ENV载入数据
- `LoadData(dataSource ...interface{}) (err error)` 从struct或map加载数据
- `LoadFlags(keys []string) (err error)` 从命令行参数载入数据
- `LoadExists(sourceFiles ...string) (err error)` 从存在的配置文件里加载数据，会忽略不存在的文件
- `LoadFiles(sourceFiles ...string) (err error)` 从给定的配置文件里加载数据，有文件不存在则会panic
- `LoadRemote(format, url string) (err error)` 从远程 URL 加载配置数据
- `LoadSources(format string, src []byte, more ...[]byte) (err error)` 从给定格式的字节数据加载配置
- `LoadStrings(format string, str string, more ...string) (err error)` 从给定格式的字符串配置里加载配置数据

### 获取值

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

### 设置值

- `Set(key string, val interface{}, setByPath ...bool) (err error)`

### 有用的方法

- `Getenv(name string, defVal ...string) (val string)`
- `AddDriver(driver Driver)`
- `Data() map[string]interface{}`
- `Exists(key string, findByPath ...bool) bool`
- `DumpTo(out io.Writer, format string) (n int64, err error)`
- `SetData(data map[string]interface{})` 设置数据以覆盖 `Config.Data`
- `BindStruct(key string, dst interface{}) error`

## 单元测试

```bash
go test -cover
// contains all sub-folder
go test -cover ./...
```

## Gookit 工具包

- [gookit/ini](https://github.com/gookit/ini) INI配置读取管理，支持多文件加载，数据覆盖合并, 解析ENV变量, 解析变量引用
- [gookit/rux](https://github.com/gookit/rux) Simple and fast request router for golang HTTP 
- [gookit/gcli](https://github.com/gookit/gcli) Go的命令行应用，工具库，运行CLI命令，支持命令行色彩，用户交互，进度显示，数据格式化显示
- [gookit/event](https://github.com/gookit/event) Go实现的轻量级的事件管理、调度程序库, 支持设置监听器的优先级, 支持对一组事件进行监听
- [gookit/cache](https://github.com/gookit/cache) 通用的缓存使用包装库，通过包装各种常用的驱动，来提供统一的使用API
- [gookit/config](https://github.com/gookit/config) Go应用配置管理，支持多种格式（JSON, YAML, TOML, INI, HCL, ENV, Flags），多文件加载，远程文件加载，数据合并
- [gookit/color](https://github.com/gookit/color) CLI 控制台颜色渲染工具库, 拥有简洁的使用API，支持16色，256色，RGB色彩渲染输出
- [gookit/filter](https://github.com/gookit/filter) 提供对Golang数据的过滤，净化，转换
- [gookit/validate](https://github.com/gookit/validate) Go通用的数据验证与过滤库，使用简单，内置大部分常用验证、过滤器
- [gookit/goutil](https://github.com/gookit/goutil) Go 的一些工具函数，格式化，特殊处理，常用信息获取等
- 更多请查看 https://github.com/gookit

## 相关包

- Ini 解析 [gookit/ini/parser](https://github.com/gookit/ini/tree/master/parser)
- Yaml 解析 [go-yaml](https://github.com/go-yaml/yaml)
- Toml 解析 [go toml](https://github.com/BurntSushi/toml)
- 数据合并 [mergo](https://github.com/imdario/mergo)
- 映射数据到结构体 [mapstructure](https://github.com/mitchellh/mapstructure)

## License

**MIT**
