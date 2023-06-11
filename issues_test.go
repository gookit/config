package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/ini"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/testutil/assert"
)

// https://github.com/gookit/config/issues/37
func TestIssues_37(t *testing.T) {
	is := assert.New(t)

	c := config.New("test")
	c.AddDriver(yaml.Driver)

	err := c.LoadStrings(config.JSON, `{
    "lang": {
        "allowed": {
            "en": "ddd"
        }
    }
}
`)
	is.NoErr(err)
	dump.Println(c.Data())

	// update yaml pkg to goccy/go-yaml
	is.NotPanics(func() {
		_ = c.LoadStrings(config.Yaml, `
lang:
  allowed:
    en: "666"
`)
	})
}

// https://github.com/gookit/config/issues/37
func TestIssues37_yaml_v3(t *testing.T) {
	is := assert.New(t)

	c := config.New("test")
	c.AddDriver(yamlv3.Driver)

	err := c.LoadStrings(config.JSON, `
{
    "lang": {
        "allowed": {
            "en": "ddd"
        }
    }
}
`)
	is.NoErr(err)
	dump.Println(c.Data())

	err = c.LoadStrings(config.Yaml, `
lang:
  newKey: hhh
  allowed:
    en: "666"
`)
	is.NoErr(err)
	dump.Println(c.Data())
}

// BindStruct doesn't seem to work with env var substitution
// https://github.com/gookit/config/issues/46
func TestIssues_46(t *testing.T) {
	is := assert.New(t)

	c := config.New("test").WithOptions(config.ParseEnv)
	err := c.LoadStrings(config.JSON, `
{
  "http": {
    "port": "${HTTP_PORT|8080}"
  }
}
`)

	is.NoErr(err)
	dump.Println(c.Data())

	val, _ := c.GetValue("http")
	mp := val.(map[string]any)
	dump.Println(mp)
	is.Eq("${HTTP_PORT|8080}", mp["port"])

	smp := c.StringMap("http")
	dump.Println(smp)
	is.Contains(smp, "port")
	is.Eq("8080", smp["port"])

	type Http struct {
		Port int
	}

	h := &Http{}
	err = c.BindStruct("http", h)
	is.NoErr(err)
	dump.Println(h)
	is.Eq(8080, h.Port)

	testutil.MockEnvValue("HTTP_PORT", "19090", func(_ string) {
		h := &Http{}
		err = c.BindStruct("http", h)
		is.NoErr(err)
		dump.Println(h)
		is.Eq(19090, h.Port)
	})
}

// https://github.com/gookit/config/issues/46
func TestIssues_59(t *testing.T) {
	is := assert.New(t)

	c := config.NewWithOptions("test", config.ParseEnv)
	c.AddDriver(ini.Driver)

	err := c.LoadFiles("testdata/ini_base.ini")
	is.NoErr(err)
	dump.Println(c.Data())

	dumpfile := "testdata/issues59.ini"
	out, err := fsutil.OpenTruncFile(dumpfile, 0666)
	is.NoErr(err)
	_, err = c.DumpTo(out, config.Ini)
	is.NoErr(err)

	str := string(fsutil.MustReadFile(dumpfile))
	is.Contains(str, "name = app")
	is.Contains(str, "key1 = val1")
}

// https://github.com/gookit/config/issues/70
func TestIssues_70(t *testing.T) {
	c := config.New("test")

	err := c.LoadStrings(config.JSON, `{
  "parent": {
    "child": "Test Var"  
  }
}`)

	assert.NoErr(t, err)
	assert.Eq(t, "Test Var", c.String("parent.child"))
	dump.P(c.Data())

	// cannot this.
	err = c.Set("parent.child.grandChild", "New Val")
	assert.Err(t, err)

	err = c.Set("parent.child", map[string]any{
		"grandChild": "New Val",
	})
	assert.NoErr(t, err)
	assert.Eq(t, map[string]any{
		"grandChild": "New Val",
	}, c.Get("parent.child"))

	dump.P(c.Data())
}

// https://github.com/gookit/config/issues/76
func TestIssues_76(t *testing.T) {
	is := assert.New(t)
	c := config.New("test")

	err := c.LoadStrings(config.JSON, `
{
    "lang": {
        "allowed": {
            "en": "ddd"
        }
    },
	"key0": 234
}
`)
	is.NoErr(err)

	ss := c.Strings("key0")
	is.Empty(ss)

	lastErr := c.Error()
	is.Err(lastErr)
	is.Eq("value cannot be convert to []string, key is 'key0'", lastErr.Error())
	is.NoErr(c.Error())
}

// https://github.com/gookit/config/issues/81
func TestIssues_81(t *testing.T) {
	is := assert.New(t)
	c := config.New("test").WithOptions(config.ParseTime, func(options *config.Options) {
		options.DecoderConfig.TagName = "json"
	})

	err := c.LoadStrings(config.JSON, `
{
	"key0": "abc",
	"age": 12,
	"connTime": "10s",
	"idleTime": "1m"
}
`)
	is.NoErr(err)

	type Options struct {
		ConnTime time.Duration `json:"connTime"`
		IdleTime time.Duration `json:"idleTime"`
	}

	opt := &Options{}
	err = c.BindStruct("", opt)

	is.NoErr(err)
	is.Eq("10s", c.String("connTime"))
	wantTm, err := time.ParseDuration("10s")
	is.NoErr(err)
	is.Eq(wantTm, opt.ConnTime)

	is.Eq("1m", c.String("idleTime"))
	wantTm, err = time.ParseDuration("1m")
	is.NoErr(err)
	is.Eq(wantTm, opt.IdleTime)
}

// https://github.com/gookit/config/issues/94
func TestIssues_94(t *testing.T) {
	is := assert.New(t)
	// add option: config.ParseDefault
	c := config.New("test").WithOptions(config.ParseDefault)

	// only set name
	c.SetData(map[string]any{
		"name": "inhere",
	})

	// age load from default tag
	type User struct {
		Age  int `json:"age" default:"30"`
		Name string
		Tags []int
	}

	user := &User{}
	is.NoErr(c.Decode(user))
	dump.Println(user)
	is.Eq("inhere", user.Name)
	is.Eq(30, user.Age)

	// field use ptr
	type User1 struct {
		Age  *int `json:"age" default:"30"`
		Name string
		Tags []int
	}

	u1 := &User1{}
	is.NoErr(c.Decode(u1))
	dump.Println(u1)
	is.Eq("inhere", u1.Name)
	is.Eq(30, *u1.Age)
}

// https://github.com/gookit/config/issues/96
func TestIssues_96(t *testing.T) {
	is := assert.New(t)
	c := config.New("test")

	err := c.Set("parent.child[0]", "Test1")
	is.NoErr(err)
	err = c.Set("parent.child[1]", "Test2")
	is.NoErr(err)

	dump.Println(c.Data())
	is.NotEmpty(c.Data())
	is.Eq([]string{"Test1", "Test2"}, c.Get("parent.child"))
}

// https://github.com/gookit/config/issues/114
func TestIssues_114(t *testing.T) {
	c := config.NewWithOptions("test",
		config.ParseDefault,
		config.ParseEnv,
		config.Readonly,
	)

	type conf struct {
		Name  string   `mapstructure:"name" default:"${NAME | Bob}"`
		Value []string `mapstructure:"value" default:"${VAL | val1}"`
	}

	err := c.LoadExists("")
	assert.NoErr(t, err)

	var cc conf
	err = c.Decode(&cc)
	assert.NoErr(t, err)

	assert.Eq(t, "Bob", cc.Name)
	assert.Eq(t, []string{"val1"}, cc.Value)
	// dump.Println(cc)
}

// https://github.com/gookit/config/issues/139
func TestIssues_139(t *testing.T) {
	c := config.New("issues_139", config.ParseEnv)
	c.AddDriver(ini.Driver)
	c.AddAlias("ini", "conf")

	err := c.LoadFiles("testdata/issues_139.conf")
	assert.NoErr(t, err)

	assert.Eq(t, "app", c.String("name"))
	assert.Eq(t, "defValue", c.String("envKey1"))
}

// https://github.com/gookit/config/issues/141
func TestIssues_141(t *testing.T) {
	type Logger struct {
		Name     string `json:"name"`
		LogFile  string `json:"logFile"`
		MaxSize  int    `json:"maxSize" default:"1024"` // MB
		MaxDays  int    `json:"maxDays" default:"7"`
		Compress bool   `json:"compress" default:"true"`
	}

	type LogConfig struct {
		Loggers []*Logger `default:""` // mark for parse default
	}

	c := config.New("issues_141", config.ParseDefault)
	err := c.LoadStrings(config.JSON, `
{
	"loggers": [
		{
			"name": "error",
			"logFile": "logs/error.log"
		},
		{	
			"name": "request",
			"logFile": "logs/request.log",
			"maxSize": 2048,
			"maxDays": 30,
			"compress": false
		}
	]
}
`)

	assert.NoErr(t, err)

	opt := &LogConfig{}
	err = c.Decode(opt)
	dump.Println(opt)
	assert.NoErr(t, err)
	assert.Eq(t, 2, len(opt.Loggers))
	assert.Eq(t, 1024, opt.Loggers[0].MaxSize)
	assert.Eq(t, 7, opt.Loggers[0].MaxDays)
	assert.Eq(t, true, opt.Loggers[0].Compress)

	assert.Eq(t, 2048, opt.Loggers[1].MaxSize)
	assert.Eq(t, 30, opt.Loggers[1].MaxDays)
	assert.Eq(t, false, opt.Loggers[1].Compress)
}

// https://github.com/gookit/config/issues/146
func TestIssues_146(t *testing.T) {
	c := config.NewWithOptions("test",
		config.ParseDefault,
		config.ParseEnv,
		config.ParseTime,
	)

	type conf struct {
		Env        time.Duration
		DefaultEnv time.Duration
		NoEnv      time.Duration
	}

	err := os.Setenv("ENV", "5s")
	assert.NoError(t, err)

	err = c.LoadStrings(config.JSON, `{
		"env": "${ENV}",
		"defaultEnv": "${DEFAULT_ENV| 10s}",
		"noEnv": "15s"
	}`)
	assert.NoErr(t, err)

	var cc conf
	err = c.Decode(&cc)
	assert.NoErr(t, err)

	assert.Eq(t, 5*time.Second, cc.Env)
	assert.Eq(t, 10*time.Second, cc.DefaultEnv)
	assert.Eq(t, 15*time.Second, cc.NoEnv)
}
