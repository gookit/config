package ini

import (
	"fmt"
	"testing"

	"github.com/gookit/config/v2"
	"github.com/stretchr/testify/assert"
)

func Example() {
	config.WithOptions(config.ParseEnv)

	// add Decoder and Encoder
	config.AddDriver(Driver)
	// Or
	// config.SetEncoder(config.Ini, ini.Encoder)

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

	name := config.String("name")
	fmt.Printf("get string\n - val: %v\n", name)

	// NOTICE: ini is not support array

	map1 := config.StringMap("map1")
	fmt.Printf("get map\n - val: %#v\n", map1)

	val0 := config.String("map1.key")
	fmt.Printf("get sub-value by path 'map.key'\n - val: %v\n", val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", config.String("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", config.String("envKey1", ""))

	// set value
	_ = config.Set("name", "new name")
	name = config.String("name")
	fmt.Printf("set string\n - val: %v\n", name)
}

func TestDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("ini", Driver.Name())
	// st.IsType(new(Encoder), JSONDriver.GetEncoder())

	c := config.NewEmpty("test")
	st.False(c.HasDecoder(config.Ini))

	c.AddDriver(Driver)
	st.True(c.HasDecoder(config.Ini))
	st.True(c.HasEncoder(config.Ini))

	_, err := Encoder(map[string]interface{}{"k": "v"})
	st.Nil(err)

	_, err = Encoder("invalid")
	st.Error(err)
}
