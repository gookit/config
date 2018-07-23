package toml

import (
	"fmt"
	"github.com/gookit/config"
	"testing"
	"github.com/stretchr/testify/assert"
)

var tomlStr = `
title = "TOML Example"
name = "app"

envKey = "${SHELL}"
envKey1 = "${NotExist|defValue}"

arr1 = [
  "alpha",
  "omega"
]

[map1]
name = "inhere"
org = "GitHub"
`

func Example() {
	config.WithOptions(config.ParseEnv)

	// add Decoder and Encoder
	config.AddDriver(Driver)

	err := config.LoadFiles("testdata/toml_base.toml")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", Data())

	// load more files
	err = config.LoadFiles("testdata/toml_other.toml")
	// can also
	// config.LoadFiles("testdata/toml_base.toml", "testdata/toml_other.toml")
	if err != nil {
		panic(err)
	}

	// load from string
	config.LoadSources(config.Toml, []byte(tomlStr))

	// fmt.Printf("config data: \n %#v\n", Data())
	fmt.Print("get config example:\n")

	name, ok := config.String("name")
	fmt.Printf("get string\n - ok: %v, val: %v\n", ok, name)

	arr1, ok := config.Strings("arr1")
	fmt.Printf("get array\n - ok: %v, val: %#v\n", ok, arr1)

	val0, ok := config.String("arr1.0")
	fmt.Printf("get sub-value by path 'arr.index'\n - ok: %v, val: %v\n", ok, val0)

	map1, ok := config.StringMap("map1")
	fmt.Printf("get map\n - ok: %v, val: %#v\n", ok, map1)

	val0, ok = config.String("map1.name")
	fmt.Printf("get sub-value by path 'map.key'\n - ok: %v, val: %v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.DefString("envKey1", ""))

	// Out:
	// get config example:
	// get string
	// - ok: true, val: app2
	// get array
	// - ok: true, val: []string{"alpha", "omega"}
	// get sub-value by path 'arr.index'
	// - ok: true, val: alpha
	// get map
	// - ok: true, val: map[string]string{"name":"inhere", "org":"GitHub"}
	// get sub-value by path 'map.key'
	// - ok: true, val: inhere
	// get env 'envKey' val: /bin/zsh
	// get env 'envKey1' val: defValue
}

func TestDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("toml", Driver.Name())

	c := config.NewEmpty("test")
	st.False(c.HasDecoder(config.Toml))

	c.AddDriver(Driver)
	st.True(c.HasDecoder(config.Toml))
	st.True(c.HasEncoder(config.Toml))

	tg := new(map[string]interface{})
	err := Decoder([]byte("invalid"), tg)
	st.Error(err)

	_, err = Encoder("invalid")
	st.Error(err)

	out, err := Encoder(map[string]interface{}{"k":"v"})
	st.Nil(err)
	st.Contains(string(out), `k = "v"`)
}
