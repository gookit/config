package config

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
        "key2": "val2"
    },
    "arr1": [
        "val",
        "val1",
        "val2"
    ]
}`

func Example() {
	WithOptions(ParseEnv)

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

	name, ok := String("name")
	fmt.Printf("- get string\n ok: %v, val: %v\n", ok, name)

	arr1, ok := Strings("arr1")
	fmt.Printf("- get array\n ok: %v, val: %#v\n", ok, arr1)

	val0, ok := String("arr1.0")
	fmt.Printf("- get sub-value by path 'arr.index'\n ok: %v, val: %#v\n", ok, val0)

	map1, ok := StringMap("map1")
	fmt.Printf("- get map\n ok: %v, val: %#v\n", ok, map1)

	val0, ok = String("map1.key")
	fmt.Printf("- get sub-value by path 'map.key'\n ok: %v, val: %#v\n", ok, val0)

	// can parse env name(ParseEnv: true)
	fmt.Printf("get env 'envKey' val: %s\n", DefString("envKey", ""))
	fmt.Printf("get env 'envKey1' val: %s\n", DefString("envKey1", ""))

	// set value
	_ = Set("name", "new name")
	name, ok = String("name")
	fmt.Printf("- set string\n ok: %v, val: %v\n", ok, name)

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
	//  ok: true, val: app
	// - get array
	//  ok: true, val: []string{"val", "val1", "val2"}
	// - get sub-value by path 'arr.index'
	//  ok: true, val: "val"
	// - get map
	//  ok: true, val: map[string]string{"key":"val", "key1":"val1", "key2":"val2"}
	// - get sub-value by path 'map.key'
	//  ok: true, val: "val"
	// get env 'envKey' val: /bin/zsh
	// get env 'envKey1' val: defValue
	// - set string
	//  ok: true, val: new name
}

func ExampleConfig_DefBool() {
	// load from string
	err := LoadSources(JSON, []byte(jsonStr))
	if err != nil {
		panic(err)
	}

	val, ok := Bool("debug")
	fmt.Printf("get 'debug', ok: %v, val: %v\n", ok, val)
	val1 := DefBool("debug", false)
	fmt.Printf("get 'debug' with default, val: %v\n", val1)

	// Output:
	// get 'debug', ok: true, val: true
	// get 'debug' with default, val: true
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

func TestLoad(t *testing.T) {
	is := assert.New(t)

	c := New("test")
	err := c.LoadExists("testdata/json_base.json", "not-exist.json")
	is.Nil(err)

	c.ClearAll()

	// load map
	err = c.LoadData(map[string]interface{}{
		"name":    "inhere",
		"age":     28,
		"working": true,
		"tags":    []string{"a", "b"},
		"info":    map[string]string{"k1": "a", "k2": "b"},
	})

	is.NotEmpty(c.Data())
	is.Nil(err)

	// LoadData
	err = c.LoadData("invalid")
	is.Error(err)

	is.Panics(func() {
		c.WithOptions(ParseEnv)
	})

	err = c.LoadStrings(JSON, `{"name": "inhere"}`, jsonStr)
	is.Nil(err)

	// LoadSources
	err = c.LoadSources(JSON, []byte(`{"name": "inhere"}`), []byte(jsonStr))
	is.Nil(err)

	err = c.LoadSources(JSON, []byte(`invalid`))
	is.Error(err)

	err = c.LoadSources(JSON, []byte(`{"name": "inhere"}`), []byte(`invalid`))
	is.Error(err)

	c = New("test")

	// LoadFiles
	err = c.LoadFiles("not-exist.json")
	is.Error(err)

	err = c.LoadFiles("testdata/json_error.json")
	is.Error(err)

	err = c.LoadExists("testdata/json_error.json")
	is.Error(err)

	// LoadStrings
	err = c.LoadStrings("invalid", jsonStr)
	is.Error(err)

	err = c.LoadStrings(JSON, "invalid")
	is.Error(err)

	err = c.LoadStrings(JSON, `{"name": "inhere"}`, "invalid")
	is.Error(err)

}

func TestConfig_LoadRemote(t *testing.T) {
	is := assert.New(t)

	// load remote config
	c := New("remote")
	url := "https://raw.githubusercontent.com/gookit/config/master/testdata/json_base.json"
	err := c.LoadRemote(JSON, url)
	is.Nil(err)
	is.Equal("123", c.DefString("age", ""))

	is.Len(c.LoadedFiles(), 1)
	is.Equal(url, c.LoadedFiles()[0])

	// load invalid remote data
	url1 := "https://raw.githubusercontent.com/gookit/config/master/testdata/json_error.json"
	err = c.LoadRemote(JSON, url1)
	is.Error(err)

	// load not exist
	url2 := "https://raw.githubusercontent.com/gookit/config/master/testdata/not-exist.txt"
	err = c.LoadRemote(JSON, url2)
	is.Error(err)

	// invalid remote url
	url3 := "invalid-url"
	err = c.LoadRemote(JSON, url3)
	is.Error(err)
}

func TestConfig_LoadFlags(t *testing.T) {
	is := assert.New(t)

	// load flag info
	c := New("flag")
	bakArgs := os.Args
	os.Args = []string{
		"./cliapp",
		"--name", "my-app",
		"--env", "dev",
		"--debug", "true",
	}

	err := c.LoadFlags([]string{"name", "env", "debug"})
	is.Nil(err)
	is.Equal("my-app", c.DefString("name", ""))
	is.Equal("dev", c.DefString("env", ""))
	is.True(c.DefBool("debug", false))

	fmt.Printf("%#v\n", c.Data())

	os.Args = bakArgs
}

func TestJSONDriver(t *testing.T) {
	st := assert.New(t)

	st.Equal("json", JSONDriver.Name())

	// empty
	c := NewEmpty("test")
	st.False(c.HasDecoder(JSON))

	c.AddDriver(JSONDriver)
	st.True(c.HasDecoder(JSON))
	st.True(c.HasEncoder(JSON))
}

func TestDriver(t *testing.T) {
	st := assert.New(t)

	c := Default()
	st.True(c.HasDecoder(JSON))
	st.True(c.HasEncoder(JSON))

	c.DelDriver(JSON)
	st.False(c.HasDecoder(JSON))
	st.False(c.HasEncoder(JSON))

	AddDriver(JSONDriver)
	st.True(c.HasDecoder(JSON))
	st.True(c.HasEncoder(JSON))

	c.DelDriver(JSON)
	c.SetDecoders(map[string]Decoder{JSON: JSONDecoder})
	c.SetEncoders(map[string]Encoder{JSON: JSONEncoder})
	st.True(c.HasDecoder(JSON))
	st.True(c.HasEncoder(JSON))
}

func TestOptions(t *testing.T) {
	st := assert.New(t)

	// options: ParseEnv
	c := New("test")
	c.WithOptions(ParseEnv)

	st.True(c.Options().ParseEnv)

	err := c.LoadStrings(JSON, jsonStr)
	st.Nil(err)

	str, ok := c.String("name")
	st.True(ok)
	st.Equal("app", str)

	str, ok = c.String("envKey")
	st.True(ok)
	st.NotContains(str, "${")

	str, ok = c.String("invalidEnvKey")
	st.True(ok)
	st.Contains(str, "${")

	str, ok = c.String("envKey1")
	st.True(ok)
	st.NotContains(str, "${")
	st.Equal("defValue", str)

	// options: Readonly
	c = New("test")
	c.WithOptions(Readonly)

	st.True(c.Options().Readonly)

	err = c.LoadStrings(JSON, jsonStr)
	st.Nil(err)

	str, ok = c.String("name")
	st.True(ok)
	st.Equal("app", str)

	err = c.Set("name", "new app")
	st.Error(err)
}

func TestExport(t *testing.T) {
	at := assert.New(t)

	c := New("test")

	str := c.ToJSON()
	at.Equal("", str)

	err := c.LoadStrings(JSON, jsonStr)
	at.Nil(err)

	str = c.ToJSON()
	at.Contains(str, `"name":"app"`)

	buf := &bytes.Buffer{}
	_, err = c.WriteTo(buf)
	at.Nil(err)

	buf = &bytes.Buffer{}

	_, err = c.DumpTo(buf, "invalid")
	at.Error(err)

	_, err = c.DumpTo(buf, Yml)
	at.Error(err)

	_, err = c.DumpTo(buf, JSON)
	at.Nil(err)
}

func TestStripJSONComments(t *testing.T) {
	is := assert.New(t)

	str := StripJSONComments(`{"name":"app"}`)
	is.Equal(`{"name":"app"}`, str)

	givens := []string{
		// single line comments
		`{
"name":"app" // comments
}`,
		`{
// comments
"name":"app" 
}`,
		`{"name":"app"} // comments
`,
		// multi line comments
		`{"name":"app"} /* comments */
`,
		`/* comments */
{"name":"app"}`,
		`/* 
comments 
*/
{"name":"app"}`,
		`/** 
comments 
*/
{"name":"app"}`,
		`/** 
comments 
**/
{"name":"app"}`,
		`/** 
* comments 
**/
{"name":"app"}`,
		`/** 
/* comments 
**/
{"name":"app"}`,
		`/** 
/* comments *
**/
{"name":"app"}`,
		`{"name": /*comments*/"app"}`,
		`{/*comments*/"name": "app"}`,
	}
	wants := []string{
		`{"name":"app"}`,
		`{"name":"app"}`,
		`{"name":"app"}`,
		// multi line comments
		`{"name":"app"}`,
		`{"name":"app"}`,
		`{"name":"app"}`,
		`{"name":"app"}`,
		`{"name":"app"}`,
		`{"name":"app"}`,
		`{"name":"app"}`,
		`{"name":"app"}`,
		`{"name": "app"}`,
		`{"name": "app"}`,
	}

	for i, s := range givens {
		is.Equal(wants[i], StripJSONComments(s))
	}

	str = StripJSONComments(`{"name":"app"} // comments`)
	is.Equal(`{"name":"app"}`, str)

	// fix https://github.com/gookit/config/issues/2
	str = StripJSONComments(`{"name":"http://abc.com"} // comments`)
	is.Equal(`{"name":"http://abc.com"}`, str)

	str = StripJSONComments(`{
"address": [
	"http://192.168.1.XXX:2379"
]
} // comments`)
	is.Equal(`{"address":["http://192.168.1.XXX:2379"]}`, str)

	s := `{"name":"http://abc.com"} // comments`
	s = StripJSONComments(s)
	assert.Equal(t, `{"name":"http://abc.com"}`, s)

	s = `
{// comments
    "name": "app", // comments
/*comments*/
    "debug": false,
    "baseKey": "value", // comments
	/* comments */
    "age": 123,
    "envKey1": "${NotExist|defValue}",
    "map1": { // comments
        "key": "val",
        "key1": "val1",
        "key2": "val2"
    },
    "arr1": [ // comments
        "val",
        "val1", // comments
		/* comments */
        "val2",
		"http://a.com"
    ],
	/* 
		comments 
*/
    "lang": {
		/** 
 		 * comments 
 		 */
        "dir": "res/lang",
        "allowed": {
            "en": "val",
            "zh-CN": "val2"
        }
    }
}`
	s = StripJSONComments(s)
	assert.Equal(
		t,
		`{"name":"app","debug":false,"baseKey":"value","age":123,"envKey1":"${NotExist|defValue}","map1":{"key":"val","key1":"val1","key2":"val2"},"arr1":["val","val1","val2","http://a.com"],"lang":{"dir":"res/lang","allowed":{"en":"val","zh-CN":"val2"}}}`,
		s,
	)
}
