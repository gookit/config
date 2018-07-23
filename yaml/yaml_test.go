package yaml

import (
	"bytes"
	"fmt"
	"github.com/gookit/config"
	"testing"
	"github.com/stretchr/testify/assert"
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
	config.WithOptions(config.ParseEnv)

	// add yaml decoder
	// only add decoder
	// config.SetDecoder(config.Yaml, Decoder)
	// Or
	config.AddDriver(Driver)

	err := config.LoadFiles("testdata/yml_other.yml")
	if err != nil {
		panic(err)
	}

	// load from string
	config.LoadSources(config.Yaml, []byte(yamlStr))

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

	// Out:
	// get config example:
	// get string
	// - ok: true, val: app2
	// get array
	// - ok: true, val: []string{"val1", "val21"}
	// get sub-value by path 'arr.index'
	// - ok: true, val: "val1"
	// get map
	// ok: true, val: map[string]string{"key":"val2", "key2":"val20"}
	// get sub-value by path 'map.key'
	// - ok: true, val: "val2"
	// get env 'envKey' val: /bin/zsh
	// get env 'envKey1' val: defValue
}

func Example_exportConfig() {
	// Notice: before dump please set driver encoder
	config.SetEncoder(config.Yaml, Encoder)

	buf := new(bytes.Buffer)
	_, err := config.DumpTo(buf, config.Yaml)
	if err != nil {
		panic(err)
	}

	fmt.Printf("export config:\n%s", buf.String())

	// Out:
	// arr1:
	// 	- val1
	// 	- val21
	// baseKey: value2
	// debug: false
	// ... ...
}

func TestDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("yaml", Driver.Name())
	// st.IsType(new(Encoder), JsonDriver.GetEncoder())

	c := config.NewEmpty("test")
	st.False(c.HasDecoder(config.Yaml))
	c.AddDriver(Driver)
	st.True(c.HasDecoder(config.Yaml))
	st.True(c.HasEncoder(config.Yaml))
}

