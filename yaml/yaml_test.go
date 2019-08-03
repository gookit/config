package yaml

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/config/v2"
	"github.com/gookit/goutil/testutil"
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
	_ = config.LoadSources(config.Yaml, []byte(yamlStr))

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

	// Out:
	// get config example:
	// get string
	// - val: app2
	// get array
	// - val: []string{"val1", "val21"}
	// get sub-value by path 'arr.index'
	// - val: "val1"
	// get map
	// val: map[string]string{"key":"val2", "key2":"val20"}
	// get sub-value by path 'map.key'
	// - val: "val2"
	// get env 'envKey' val: /bin/zsh
	// get env 'envKey1' val: defValue
}

func Example_exportConfig() {
	// Notice: before dump please set driver encoder
	config.AddDriver(Driver)

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
	// st.IsType(new(Encoder), JSONDriver.GetEncoder())

	c := config.NewEmpty("test")
	st.False(c.HasDecoder(config.Yaml))
	c.AddDriver(Driver)
	st.True(c.HasDecoder(config.Yaml))
	st.True(c.HasEncoder(config.Yaml))
}

// Support "=", ":", "." characters for default values
// see https://github.com/gookit/config/issues/9
func TestIssue2(t *testing.T) {
	ris := assert.New(t)

	c := config.NewEmpty("test")
	c.AddDriver(Driver)
	c.WithOptions(config.ParseEnv)

	err := c.LoadStrings(config.Yaml, `
command: ${APP_COMMAND|app:run}
`)
	ris.NoError(err)
	testutil.MockEnvValue("APP_COMMAND", "new val", func(nv string) {
		ris.Equal("new val", nv)
		ris.Equal("new val", c.String("command"))
	})

	ris.Equal("", config.Getenv("APP_COMMAND"))
	ris.Equal("app:run", c.String("command"))

	c.ClearAll()
	err = c.LoadStrings(config.Yaml, `
command: ${ APP_COMMAND | app:run }
`)
	ris.NoError(err)
	testutil.MockEnvValue("APP_COMMAND", "new val", func(nv string) {
		ris.Equal("new val", nv)
		ris.Equal("new val", c.String("command"))
	})
	ris.Equal("", config.Getenv("APP_COMMAND"))
	ris.Equal("app:run", c.String("command"))
}
