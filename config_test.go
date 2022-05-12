package config

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/gookit/goutil/testutil"
	"github.com/stretchr/testify/assert"
)

var jsonStr = `{
    "name": "app",
    "debug": true,
    "baseKey": "value",
    "age": 123,
    "envKey": "${SHELL}",
    "envKey1": "${NotExist|defValue}",
    "invalidEnvKey": "${noClose",
    "map1": {
        "key": "val",
        "key1": "val1",
        "key2": "val2",
        "key3": "${SHELL}"
    },
    "arr1": [
        "val",
        "val1",
        "val2"
    ]
}`

func Example() {
	// WithOptions(ParseEnv)

	// add Decoder and Encoder
	// use yaml github.com/gookit/config/yaml
	// AddDriver(Yaml, yaml.Driver)
	// use toml github.com/gookit/config/toml
	// AddDriver(Toml, toml.Driver)
	// use toml github.com/gookit/config/hcl
	// AddDriver(Hcl, hcl.Driver)
	// Or
	// config.DecoderEncoder(config.JSON, yaml.Decoder, yaml.Encoder)

	err := LoadFiles("testdata/json_base.json")
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", Data())

	err = LoadFiles("testdata/json_other.json")
	// LoadFiles("testdata/json_base.json", "testdata/json_other.json")
	if err != nil {
		panic(err)
	}

	// load from string
	err = LoadSources(JSON, []byte(jsonStr))
	if err != nil {
		panic(err)
	}

	// fmt.Printf("config data: \n %#v\n", Data())
	fmt.Print("get config example:\n")

	name := String("name")
	fmt.Printf("- get string\n val: %v\n", name)

	arr1 := Strings("arr1")
	fmt.Printf("- get array\n val: %#v\n", arr1)

	val0 := String("arr1.0")
	fmt.Printf("- get sub-value by path 'arr.index'\n val: %#v\n", val0)

	map1 := StringMap("map1")
	fmt.Printf("- get map\n val: %#v\n", map1)

	val0 = String("map1.key")
	fmt.Printf("- get sub-value by path 'map.key'\n val: %#v\n", val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", String("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", String("envKey1", ""))

	// set value
	_ = Set("name", "new name")
	name = String("name")
	fmt.Printf("- set string\n val: %v\n", name)

	// if you want export config data
	// buf := new(bytes.Buffer)
	// _, err = config.DumpTo(buf, config.JSON)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("export config:\n%s", buf.String())

	// Out:
	// get config example:
	// - get string
	//  val: app
	// - get array
	//  val: []string{"val", "val1", "val2"}
	// - get sub-value by path 'arr.index'
	//  val: "val"
	// - get map
	//  val: map[string]string{"key":"val", "key1":"val1", "key2":"val2"}
	// - get sub-value by path 'map.key'
	//  val: "val"
	// get env 'envKey' val: /bin/zsh
	// get env 'envKey1' val: defValue
	// - set string
	//  val: new name
}

func Example_exportConfig() {
	// Notice: before dump please set driver encoder
	// SetEncoder(Yaml, yaml.Encoder)

	ClearAll()
	// load from string
	err := LoadStrings(JSON, `{
"name": "app",
"age": 34
}`)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	_, err = DumpTo(buf, JSON)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", buf.String())

	// Output:
	// {"age":34,"name":"app"}
}

func BenchmarkGet(b *testing.B) {
	err := LoadStrings(JSON, jsonStr)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Get("name")
	}
}

func TestBasic(t *testing.T) {
	st := assert.New(t)

	ClearAll()
	c := Default()
	st.True(c.HasDecoder(JSON))
	st.True(c.HasEncoder(JSON))
	st.Equal("default", c.Name())
	st.NoError(c.Error())

	c = NewWithOptions("test", Readonly)
	opts := c.Options()
	st.True(opts.Readonly)
	st.Equal(JSON, opts.DumpFormat)
	st.Equal(JSON, opts.ReadFormat)
}

func TestGetEnv(t *testing.T) {
	testutil.MockEnvValues(map[string]string{
		"APP_NAME":  "config",
		"APP_DEBUG": "true",
	}, func() {
		assert.Equal(t, "config", Getenv("APP_NAME"))
		assert.Equal(t, "true", Getenv("APP_DEBUG"))
		assert.Equal(t, "defVal", GetEnv("not-exsit", "defVal"))
	})

	assert.Equal(t, Yaml, fixFormat("yml"))
	assert.Equal(t, Hcl, fixFormat("conf"))
	assert.Equal(t, Ini, fixFormat("inc"))
}

func TestSetDecoderEncoder(t *testing.T) {
	is := assert.New(t)

	c := Default()
	c.ClearAll()

	is.True(c.HasDecoder(JSON))
	is.True(c.HasEncoder(JSON))

	c.DelDriver(JSON)

	is.False(c.HasDecoder(JSON))
	is.False(c.HasEncoder(JSON))

	SetDecoder(JSON, JSONDecoder)
	SetEncoder(JSON, JSONEncoder)

	is.True(c.HasDecoder(JSON))
	is.True(c.HasEncoder(JSON))
}

func TestDefault(t *testing.T) {
	is := assert.New(t)

	ClearAll()
	WithOptions(ParseEnv)

	is.True(GetOptions().ParseEnv)

	_ = LoadStrings(JSON, `{"name": "inhere"}`)

	buf := &bytes.Buffer{}
	_, err := WriteTo(buf)
	is.Nil(err)
}

func TestJSONDriver(t *testing.T) {
	is := assert.New(t)
	is.Equal("json", JSONDriver.Name())

	// empty
	c := NewEmpty("test")
	is.False(c.HasDecoder(JSON))

	c.AddDriver(JSONDriver)
	is.True(c.HasDecoder(JSON))
	is.True(c.HasEncoder(JSON))
	is.Len(c.DriverNames(), 1)

	is.Equal(byte('.'), c.Options().Delimiter)
	is.Equal(".", string(c.Options().Delimiter))
	c.WithOptions(func(opt *Options) {
		opt.Delimiter = 0
	})
	is.Equal(byte(0), c.Options().Delimiter)

	err := c.LoadStrings(JSON, `{"key": 1}`)
	is.NoError(err)
	is.Equal(1, c.Int("key"))

	c = NewWith("test", func(c *Config) {
		err = c.LoadData(map[string]interface{}{"key1": 2})
		is.NoError(err)
	})
	is.Equal(2, c.Int("key1"))
}

func TestDriver(t *testing.T) {
	is := assert.New(t)

	c := Default()
	is.True(c.HasDecoder(JSON))
	is.True(c.HasEncoder(JSON))

	c.DelDriver(JSON)
	is.False(c.HasDecoder(JSON))
	is.False(c.HasEncoder(JSON))

	AddDriver(JSONDriver)
	is.True(c.HasDecoder(JSON))
	is.True(c.HasEncoder(JSON))

	c.DelDriver(JSON)
	c.SetDecoders(map[string]Decoder{JSON: JSONDecoder})
	c.SetEncoders(map[string]Encoder{JSON: JSONEncoder})
	is.True(c.HasDecoder(JSON))
	is.True(c.HasEncoder(JSON))
}

func TestOptions(t *testing.T) {
	is := assert.New(t)

	// options: ParseEnv
	c := New("test")
	c.WithOptions(ParseEnv)

	is.True(c.Options().ParseEnv)

	err := c.LoadStrings(JSON, jsonStr)
	is.Nil(err)

	str := c.String("name")
	is.Equal("app", str)

	// test: parse env name
	shell := os.Getenv("SHELL")
	// ensure env var is exist
	if shell == "" {
		_ = os.Setenv("SHELL", "/usr/bin/bash")
	}

	str = c.String("envKey")
	is.NotContains(str, "${")

	// revert
	if shell != "" {
		_ = os.Setenv("SHELL", shell)
	}

	str = c.String("invalidEnvKey")
	is.Contains(str, "${")

	str = c.String("envKey1")
	is.NotContains(str, "${")
	is.Equal("defValue", str)

	// options: Readonly
	c = New("test")
	c.WithOptions(Readonly)

	is.True(c.Options().Readonly)

	err = c.LoadStrings(JSON, jsonStr)
	is.Nil(err)

	str = c.String("name")
	is.Equal("app", str)

	err = c.Set("name", "new app")
	is.Error(err)
}

func TestDelimiter(t *testing.T) {
	// options: Delimiter
	is := assert.New(t)
	c := New("test")
	c.WithOptions(Delimiter(':'))
	is.Equal(byte(':'), c.Options().Delimiter)

	err := c.LoadData(map[string]interface{}{
		"top0": 1,
		"top1": map[string]int{"sub0": 2},
	})
	is.NoError(err)
	// is.Equal(1, c.Int("top0"))
	is.Equal(2, c.Int("top1:sub0"))

	// load will use defaultDelimiter
	c = NewWithOptions("test", Delimiter(0))
	is.Equal(byte(0), c.Options().Delimiter)

	err = c.LoadData(map[string]interface{}{
		"top0": 1,
		"top1": map[string]int{"sub0": 2},
	})
	is.NoError(err)
	is.Equal(2, c.Int("top1.sub0"))
}

func TestEnableCache(t *testing.T) {
	is := assert.New(t)

	c := NewWithOptions("test", EnableCache)
	err := c.LoadStrings(JSON, jsonStr)
	is.Nil(err)

	str := c.String("name")
	is.Equal("app", str)

	// re-get, from caches
	str = c.String("name")
	is.Equal("app", str)

	sArr := c.Strings("arr1")
	is.Equal("val1", sArr[1])

	// re-get, from caches
	sArr = c.Strings("arr1")
	is.Equal("val1", sArr[1])

	sMap := c.StringMap("map1")
	is.Equal("val1", sMap["key1"])
	sMap = c.StringMap("map1")
	is.Equal("val1", sMap["key1"])

	c.ClearAll()
}

func TestJSONAllowComments(t *testing.T) {
	is := assert.New(t)

	m := struct {
		N string
	}{}

	// disable clear comments
	old := JSONAllowComments
	JSONAllowComments = false
	err := JSONDecoder([]byte(`{
// comments
"n":"v"}
`), &m)
	is.Error(err)

	JSONAllowComments = old
}

func TestMapStringStringParseEnv(t *testing.T) {
	is := assert.New(t)
	c := New("test")
	c.WithOptions(ParseEnv)
	err := c.LoadStrings(JSON, jsonStr)
	is.Nil(err)

	shellVal := "/usr/bin/bash"
	testutil.MockEnvValue("SHELL", shellVal, func(_ string) {
		sMap := c.StringMap("map1")
		is.Equal(shellVal, sMap["key3"])
	})
}
