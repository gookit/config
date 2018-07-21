package json

import (
	"fmt"
	"github.com/gookit/config"
	"testing"
	"github.com/stretchr/testify/assert"
)

func Example() {
	config.WithOptions(config.WithParseEnv)

	// add Decoder and Encoder
	config.AddDriver(config.Json, Driver)
	// Or
	// config.DecoderEncoder(config.Json, json.Decoder, json.Encoder)

	err := config.LoadFiles("testdata/json_base.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())

	err = config.LoadFiles("testdata/json_other.json")
	// config.LoadFiles("testdata/json_base.json", "testdata/json_other.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config data: \n %#v\n", config.Data())
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

	// set value
	config.Set("name", "new name")
	name, ok = config.String("name")
	fmt.Printf("set string\n - ok: %v, val: %v\n", ok, name)

	// if you want export config data
	// buf := new(bytes.Buffer)
	// _, err = config.DumpTo(buf, config.Json)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("export config:\n%s", buf.String())
}

func TestDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("json", Driver.Name())
	// st.IsType(new(Encoder), Driver.GetEncoder())
}
