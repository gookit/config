package other

import (
	"fmt"
	"testing"

	"github.com/gookit/config/v2"
	"github.com/stretchr/testify/assert"
)

func TestOtherDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("other", Driver.Name())

	c := config.NewEmpty("test")
	st.False(c.HasDecoder("other"))

	c.AddDriver(Driver)
	st.True(c.HasDecoder("other"))
	st.True(c.HasEncoder("other"))

	_, err := Encoder(map[string]interface{}{"k": "v"})
	st.Nil(err)

	_, err = Encoder("invalid")
	st.Error(err)
}

func TestOtherLoader(t *testing.T) {
	config.AddDriver(Driver)

	err := config.LoadFiles("../testdata/ini_base.other")
	if err != nil {
		panic(err)
	}

	fmt.Printf("get config example:\n")

	name := config.String("name")
	fmt.Printf("get string\n - val: %v\n", name)

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
