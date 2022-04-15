package config

import (
	"os"
	"runtime"
	"testing"

	"github.com/gookit/goutil/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLoad(t *testing.T) {
	st := assert.New(t)

	ClearAll()
	err := LoadFiles("testdata/json_base.json", "testdata/json_other.json")
	st.Nil(err)

	ClearAll()
	err = LoadFilesByFormat(JSON, "testdata/json_base.json", "testdata/json_other.json")
	st.Nil(err)

	ClearAll()
	err = LoadExists("testdata/json_base.json", "not-exist.json")
	st.Nil(err)

	ClearAll()
	err = LoadExistsByFormat(JSON, "testdata/json_base.json", "not-exist.json")
	st.Nil(err)

	ClearAll()
	// load map
	err = LoadData(map[string]interface{}{
		"name":    "inhere",
		"age":     28,
		"working": true,
		"tags":    []string{"a", "b"},
		"info":    map[string]string{"k1": "a", "k2": "b"},
	})
	st.NotEmpty(Data())
	st.Nil(err)
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
	is.Equal(OnCleanData, name)

	// load map data
	err = c.LoadData(map[string]interface{}{
		"name":    "inhere",
		"age":     float64(28),
		"working": true,
		"tags":    []string{"a", "b"},
		"info":    map[string]string{"k1": "a", "k2": "b"},
	})

	is.Equal(OnLoadData, name)
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

func TestLoadRemote(t *testing.T) {
	is := assert.New(t)

	// invalid remote url
	url3 := "invalid-url"
	err := LoadRemote(JSON, url3)
	is.Error(err)

	if runtime.GOOS == "windows" {
		return
	}

	// load remote config
	c := New("remote")
	url := "https://raw.githubusercontent.com/gookit/config/master/testdata/json_base.json"
	err = c.LoadRemote(JSON, url)
	is.Nil(err)
	is.Equal("123", c.String("age", ""))

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
	is.Equal("inhere", c.String("name", ""))
	is.Equal("dev", c.String("env", ""))
	is.Equal(99, c.Int("age"))
	is.Equal(uint(12), c.Uint("var0"))
	is.Equal(uint(20), c.Uint("not-exist", uint(20)))
	is.Equal("val", c.Get("unknownTyp"))
	is.True(c.Bool("debug", false))

	// set sub key
	c = New("flag")
	_ = c.LoadStrings(JSON, jsonStr)
	os.Args = []string{
		"./binFile",
		"--map1.key", "new val",
	}
	is.Equal("val", c.String("map1.key"))
	err = c.LoadFlags([]string{"--map1.key"})
	is.NoError(err)
	is.Equal("new val", c.String("map1.key"))
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
		assert.Equal(t, "", String("test_env0"))

		LoadOSEnv([]string{"APP_NAME", "app_debug", "test_env0"}, true)

		assert.True(t, Bool("app_debug"))
		assert.Equal(t, "config", String("app_name"))
		assert.Equal(t, "val0", String("test_env0"))
		assert.Equal(t, "", String("test_env1"))
	})

	ClearAll()
}
