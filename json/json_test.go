package json

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
	// _, err = config.DumpTo(buf, config.JSON)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("export config:\n%s", buf.String())
}

func TestDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("json", Driver.Name())

	c := config.NewEmpty("test")
	st.False(c.HasDecoder(config.JSON))
	c.AddDriver(Driver)

	st.True(c.HasDecoder(config.JSON))
	st.True(c.HasEncoder(config.JSON))

	m := struct {
		N string
	}{}
	err := Decoder([]byte(`{
// comments
"n":"v"}
`), &m)
	st.Nil(err)
	st.Equal("v", m.N)

	// disable clear comments
	old := config.JSONAllowComments
	config.JSONAllowComments = false
	err = Decoder([]byte(`{
// comments
"n":"v"}
`), &m)
	st.Error(err)

	config.JSONAllowComments = old
}
