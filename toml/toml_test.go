package toml

import (
	"fmt"
	"testing"

	"github.com/gookit/config/v2"
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
	_ = config.LoadSources(config.Toml, []byte(tomlStr))

	// fmt.Printf("config data: \n %#v\n", Data())
	fmt.Print("get config example:\n")

	name := config.String("name")
	fmt.Printf("get string\n - val: %v\n", name)

	arr1 := config.Strings("arr1")
	fmt.Printf("get array\n - val: %#v\n", arr1)

	val0 := config.String("arr1.0")
	fmt.Printf("get sub-value by path 'arr.index'\n - val: %v\n", val0)

	map1 := config.StringMap("map1")
	fmt.Printf("get map\n - val: %#v\n", map1)

	val0 = config.String("map1.name")
	fmt.Printf("get sub-value by path 'map.key'\n - val: %v\n", val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.String("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.String("envKey1", ""))

	// Out:
	// get config example:
	// get string
	// - val: app2
	// get array
	// - val: []string{"alpha", "omega"}
	// get sub-value by path 'arr.index'
	// - val: alpha
	// get map
	// - val: map[string]string{"name":"inhere", "org":"GitHub"}
	// get sub-value by path 'map.key'
	// - val: inhere
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

	out, err := Encoder(map[string]interface{}{"k": "v"})
	st.Nil(err)
	st.Contains(string(out), `k = "v"`)
}
