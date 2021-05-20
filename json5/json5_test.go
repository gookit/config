package json5_test

import (
	"fmt"
	"testing"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/json5"
	"github.com/stretchr/testify/assert"
)

func Example() {
	config.WithOptions(config.ParseEnv)

	// add Decoder and Encoder
	config.AddDriver(json5.Driver)

	err := config.LoadFiles("testdata/json_base.json5")
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

	name := config.String("name")
	fmt.Printf("get string\n - val: %v\n", name)

	arr1 := config.Strings("arr1")
	fmt.Printf("get array\n - val: %#v\n", arr1)

	val0 := config.String("arr1.0")
	fmt.Printf("get sub-value by path 'arr.index'\n - val: %#v\n", val0)

	map1 := config.StringMap("map1")
	fmt.Printf("get map\n - val: %#v\n", map1)

	val0 = config.String("map1.key")
	fmt.Printf("get sub-value by path 'map.key'\n - val: %#v\n", val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.String("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.String("envKey1", ""))

	// set value
	_ = config.Set("name", "new name")
	name = config.String("name")
	fmt.Printf("set string\n - val: %v\n", name)

	// if you want export config data
	// buf := new(bytes.Buffer)
	// _, err = config.DumpTo(buf, json5.NAME)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("export config:\n%s", buf.String())
}

func TestDriver(t *testing.T) {
	is := assert.New(t)

	is.Equal(json5.NAME, json5.Driver.Name())

	c := config.NewEmpty("test")
	is.False(c.HasDecoder(json5.NAME))
	c.AddDriver(json5.Driver)

	is.True(c.HasDecoder(json5.NAME))
	is.True(c.HasEncoder(json5.NAME))

	m := struct {
		N string
	}{}
	err := json5.Decoder([]byte(`{
// comments
"n":"v"}
`), &m)
	is.Nil(err)
	is.Equal("v", m.N)
}
