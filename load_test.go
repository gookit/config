package config

import (
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/testutil/assert"
)

func TestDefaultLoad(t *testing.T) {
	is := assert.New(t)

	ClearAll()
	err := LoadFiles("testdata/json_base.json", "testdata/json_other.json")
	is.Nil(err)

	ClearAll()
	err = LoadFilesByFormat(JSON, "testdata/json_base.json", "testdata/json_other.json")
	is.Nil(err)

	ClearAll()
	err = LoadExists("testdata/json_base.json", "not-exist.json")
	is.Nil(err)

	ClearAll()
	err = LoadExistsByFormat(JSON, "testdata/json_base.json", "not-exist.json")
	is.Nil(err)

	ClearAll()
	// load map
	err = LoadData(map[string]any{
		"name":    "inhere",
		"age":     28,
		"working": true,
		"tags":    []string{"a", "b"},
		"info":    map[string]string{"k1": "a", "k2": "b"},
	})
	is.NotEmpty(Data())
	is.NotEmpty(Keys())
	is.Empty(Sub("not-exist"))
	is.Nil(err)
}

func TestLoad(t *testing.T) {
	is := assert.New(t)

	var name string
	c := New("test").
		WithOptions(WithHookFunc(func(event string, c *Config) {
			name = event
		}))
	err := c.LoadExists("testdata/json_base.json", "not-exist.json")
	is.Nil(err)

	c.ClearAll()
	is.Eq(OnCleanData, name)

	// load map data
	err = c.LoadData(map[string]any{
		"name":    "inhere",
		"age":     float64(28),
		"working": true,
		"tags":    []string{"a", "b"},
		"info":    map[string]string{"k1": "a", "k2": "b"},
	}, map[string]string{"str-map": "value"})

	is.Eq(OnLoadData, name)
	is.NotEmpty(c.Data())
	is.Nil(err)
	is.Eq("value", c.String("str-map"))

	// LoadData
	err = c.LoadData("invalid")
	is.Err(err)

	is.Panics(func() {
		c.WithOptions(ParseEnv)
	})

	err = c.LoadStrings(JSON, `{"name": "inhere"}`, jsonStr)
	is.Nil(err)

	// LoadSources
	err = c.LoadSources(JSON, []byte(`{"name": "inhere"}`), []byte(jsonStr))
	is.Nil(err)

	err = c.LoadSources(JSON, []byte(`invalid`))
	is.Err(err)

	err = c.LoadSources(JSON, []byte(`{"name": "inhere"}`), []byte(`invalid`))
	is.Err(err)

	c = New("test")

	// LoadFiles
	err = c.LoadFiles("not-exist.json")
	is.Err(err)

	err = c.LoadFiles("testdata/json_error.json")
	is.Err(err)

	err = c.LoadExists("testdata/json_error.json")
	is.Err(err)

	// LoadStrings
	err = c.LoadStrings("invalid", jsonStr)
	is.Err(err)

	err = c.LoadStrings(JSON, "invalid")
	is.Err(err)

	err = c.LoadStrings(JSON, `{"name": "inhere"}`, "invalid")
	is.Err(err)
}

func TestLoadRemote(t *testing.T) {
	is := assert.New(t)

	// invalid remote url
	url3 := "invalid-url"
	err := LoadRemote(JSON, url3)
	is.Err(err)

	if runtime.GOOS == "windows" {
		t.Skip("skip test load remote on Windows")
		return
	}

	// load remote config
	c := New("remote")
	url := "https://raw.githubusercontent.com/gookit/config/master/testdata/json_base.json"
	err = c.LoadRemote(JSON, url)
	is.Nil(err)
	is.Eq("123", c.String("age", ""))

	is.Len(c.LoadedUrls(), 1)
	is.Eq(url, c.LoadedUrls()[0])

	// load invalid remote data
	url1 := "https://raw.githubusercontent.com/gookit/config/master/testdata/json_error.json"
	err = c.LoadRemote(JSON, url1)
	is.Err(err)

	// load not exist
	url2 := "https://raw.githubusercontent.com/gookit/config/master/testdata/not-exist.txt"
	err = c.LoadRemote(JSON, url2)
	is.Err(err)
}

func TestLoadFlags(t *testing.T) {
	is := assert.New(t)

	ClearAll()
	c := Default()
	bakArgs := os.Args
	// --name inhere --env dev --age 99 --debug
	os.Args = []string{
		"./binFile",
		"--env", "dev",
		"--age", "99",
		"--var0", "12",
		"--name", "inhere",
		"--unknownTyp", "val",
		"--debug",
	}

	// load flag info
	keys := []string{"name", "env", "debug:bool", "age:int", "var0:uint", "unknownTyp:notExist"}
	err := LoadFlags(keys)
	is.Nil(err)
	is.Eq("inhere", c.String("name", ""))
	is.Eq("dev", c.String("env", ""))
	is.Eq(99, c.Int("age"))
	is.Eq(uint(12), c.Uint("var0"))
	is.Eq(uint(20), c.Uint("not-exist", uint(20)))
	is.Eq("val", c.Get("unknownTyp"))
	is.True(c.Bool("debug", false))

	// set sub key
	c = New("flag")
	_ = c.LoadStrings(JSON, jsonStr)
	os.Args = []string{
		"./binFile",
		"--map1.key", "new val",
	}
	is.Eq("val", c.String("map1.key"))
	err = c.LoadFlags([]string{"--map1.key"})
	is.NoErr(err)
	is.Eq("new val", c.String("map1.key"))
	// fmt.Println(err)
	// fmt.Printf("%#v\n", c.Data())

	os.Args = bakArgs
}

func TestLoadOSEnv(t *testing.T) {
	ClearAll()

	testutil.MockEnvValues(map[string]string{
		"APP_NAME":  "config",
		"app_debug": "true",
		"test_env0": "val0",
		"TEST_ENV1": "val1",
	}, func() {
		assert.Eq(t, "", String("test_env0"))

		LoadOSEnv([]string{"APP_NAME", "app_debug", "test_env0"}, true)

		assert.True(t, Bool("app_debug"))
		assert.Eq(t, "config", String("app_name"))
		assert.Eq(t, "val0", String("test_env0"))
		assert.Eq(t, "", String("test_env1"))
	})

	ClearAll()
}

func TestLoadOSEnvs(t *testing.T) {
	ClearAll()

	testutil.MockEnvValues(map[string]string{
		"APP_NAME":  "config",
		"APP_DEBUG": "true",
		"TEST_ENV0": "val0",
		"TEST_ENV1": "val1",
	}, func() {
		assert.Eq(t, "", String("test_env0"))
		assert.Eq(t, "val0", Getenv("TEST_ENV0"))

		LoadOSEnvs(map[string]string{
			"APP_NAME":  "",
			"APP_DEBUG": "app_debug",
			"TEST_ENV0": "test0",
		})

		assert.True(t, Bool("app_debug"))
		assert.Eq(t, "config", String("app_name"))
		assert.Eq(t, "val0", String("test0"))
		assert.Eq(t, "", String("test_env1"))
	})

	ClearAll()
}

func TestLoadFromDir(t *testing.T) {
	ClearAll()
	assert.NoErr(t, LoadStrings(JSON, `{
"topKey": "a value"
}`))

	assert.NoErr(t, LoadFromDir("testdata/subdir", JSON))
	dump.P(Data())
	assert.Eq(t, "value in sub data", Get("subdata.key01"))
	assert.Eq(t, "value in task.json", Get("task.key01"))

	ClearAll()
	assert.NoErr(t, LoadFromDir("testdata/emptydir", JSON))

	// with DataKey option. see https://github.com/gookit/config/issues/173
	assert.NoErr(t, LoadFromDir("testdata/subdir", JSON, func(lo *LoadOptions) {
		lo.DataKey = "dataList"
	}))
	dump.P(Data())
	dl := Get("dataList")
	assert.NotNil(t, dl)
	assert.IsKind(t, reflect.Slice, dl)
	ClearAll()
}

func TestReloadFiles(t *testing.T) {
	ClearAll()
	c := Default()
	// no loaded files
	assert.NoErr(t, ReloadFiles())

	var eventName string
	c.WithOptions(WithHookFunc(func(event string, c *Config) {
		eventName = event
	}))

	// load files
	err := LoadFiles("testdata/json_base.json", "testdata/json_other.json")
	assert.NoErr(t, err)
	assert.Eq(t, OnLoadData, eventName)
	assert.NotEmpty(t, c.LoadedFiles())
	assert.Eq(t, "app2", c.String("name"))

	// set value
	assert.NoErr(t, c.Set("name", "new value"))
	assert.Eq(t, OnSetValue, eventName)
	assert.Eq(t, "new value", c.String("name"))

	// reload files
	assert.NoErr(t, ReloadFiles())
	assert.Eq(t, OnReloadData, eventName)

	// value is reverted
	assert.Eq(t, "app2", c.String("name"))
	ClearAll()
}
