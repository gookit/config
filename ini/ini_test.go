package ini

import (
	"fmt"
	"github.com/gookit/config"
	"testing"
	"github.com/stretchr/testify/assert"
)

func Example() {
	config.WithOptions(config.WithParseEnv)

	// add Decoder and Encoder
	config.AddDriver(config.Ini, Driver)
	// Or
	// config.DecoderEncoder(config.Ini, ini.Decoder, ini.Encoder)

	err := config.LoadFiles("testdata/ini_base.ini")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", config.Data())

	err = config.LoadFiles("testdata/ini_other.ini")
	// config.LoadFiles("testdata/ini_base.ini", "testdata/ini_other.ini")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", config.Data())
	fmt.Print("get config example:\n")

	name, ok := config.String("name")
	fmt.Printf("get string\n - ok: %v, val: %v\n", ok, name)

	// NOTICE: ini is not support array

	map1, ok := config.StringMap("map1")
	fmt.Printf("get map\n - ok: %v, val: %#v\n", ok, map1)

	val0, ok := config.String("map1.key")
	fmt.Printf("get sub-value by path 'map.key'\n - ok: %v, val: %v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.DefString("envKey1", ""))

	// set value
	config.Set("name", "new name")
	name, ok = config.String("name")
	fmt.Printf("set string\n - ok: %v, val: %v\n", ok, name)
}

func TestDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("ini", Driver.Name())
	// st.IsType(new(Encoder), JsonDriver.GetEncoder())

	c := config.NewEmpty("test")
	st.False(c.HasDecoder(config.Ini))
	st.Panics(func() {
		c.AddDriver("invalid", Driver)
	})
	c.AddDriver(config.Ini, Driver)
	st.True(c.HasDecoder(config.Ini))
	st.True(c.HasEncoder(config.Ini))
}
