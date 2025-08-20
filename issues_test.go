package config_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/nstrlabs/config"
	"github.com/nstrlabs/config/ini"
	"github.com/nstrlabs/config/yaml"
	"github.com/nstrlabs/config/yamlv3"
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
	assert.Eq(t, true, opt.Loggers[1].Compress)

	t.Run("3 elements", func(t *testing.T) {
		jsonStr := `
{
  "loggers": [
    {
      "name": "info",
      "logFile": "logs/info.log"
    },
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
`
		c := config.New("issues_141", config.ParseDefault)
		err := c.LoadStrings(config.JSON, jsonStr)
		assert.NoErr(t, err)

		opt := &LogConfig{}
		err = c.Decode(opt)
		dump.Println(opt)
	})
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

type DurationStruct struct {
	Duration time.Duration
}

// https://github.com/gookit/config/pull/151
func TestDuration(t *testing.T) {
	var (
		err error
		str string
	)

	c := config.New("test").WithOptions(config.ParseTime)
	is := assert.New(t)
	dur := DurationStruct{}

	for _, seconds := range []int{10, 90} {
		str = fmt.Sprintf(`{ "Duration": "%ds" }`, seconds)

		err = c.LoadSources(config.JSON, []byte(str))
		is.Nil(err)
		err = c.Decode(&dur)
		is.Nil(err)
		is.Equal(float64(seconds), dur.Duration.Seconds())
	}
}

// https://github.com/gookit/config/issues/156
func TestIssues_156(t *testing.T) {
	c := config.New("test", config.ParseEnv)
	c.AddDriver(yaml.Driver)

	type DbConfig struct {
		Url      string
		Type     string
		Password string
		Username string
	}

	err := c.LoadStrings(config.Yaml, `
---
datasource:
  password: ${DATABASE_PASSWORD|?} # use fixed error message
  type: postgres
  username: ${DATABASE_USERNAME|postgres}
  url: ${DATABASE_URL|?error message2}
`)
	assert.NoErr(t, err)
	// dump.Println(c.Data())
	assert.NotEmpty(t, c.Sub("datasource"))

	// will error
	dbConf := &DbConfig{}
	err = c.BindStruct("datasource", dbConf)
	assert.Err(t, err)
	assert.ErrSubMsg(t, err, "'Password' value is required for var: DATABASE_PASSWORD")
	assert.ErrSubMsg(t, err, "'Url' error message2")

	testutil.MockEnvValues(map[string]string{
		"DATABASE_PASSWORD": "1234yz56",
		"DATABASE_URL":      "localhost:5432/postgres?sslmode=disable",
	}, func() {
		dbConf := &DbConfig{}
		err = c.BindStruct("datasource", dbConf)
		assert.NoErr(t, err)
		dump.Println(dbConf)
		assert.Eq(t, "1234yz56", dbConf.Password)
		assert.Eq(t, "localhost:5432/postgres?sslmode=disable", dbConf.Url)
	})
}

// https://github.com/gookit/config/issues/162
func TestIssues_162(t *testing.T) {
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

	c := config.New("issues_162", config.ParseDefault)
	err := c.LoadStrings(config.JSON, `{}`)
	assert.NoErr(t, err)

	opt := &LogConfig{}
	err = c.Decode(opt)
	// dump.Println(opt)
	assert.Empty(t, opt.Loggers)
}

// https://github.com/gookit/goutil/issues/135
func TestGoutil_issues_135(t *testing.T) {
	// TIP: not support use JSON as default value
	testYml := `
test:
    credentials: >
        ${CREDENTIALS|{}}
    apiKey: ${API_KEY|AN_APIKEY}
    apiUri: ${API_URI|http://localhost:8888/v1/api}
`

	type Setup struct {
		Credentials string `mapstructure:"credentials"`
		ApiKey      string `mapstructure:"apiKey"`
		ApiUri      string `mapstructure:"apiUri"`
	}

	type Configuration struct {
		Details Setup `mapstructure:"test"`
	}

	c := config.New("config", config.ParseEnv).WithDriver(yamlv3.Driver)

	err := c.LoadStrings(config.Yaml, testYml)
	assert.NoErr(t, err)

	// no env values
	t.Run("no env values", func(t *testing.T) {
		st := Configuration{}
		err = c.Decode(&st)
		assert.NoErr(t, err)
		dump.Println(st)
	})

	// set value
	err = c.Set("test.credentials", `${CREDENTIALS}`)
	assert.NoErr(t, err)

	// set value(use JSON as default value)
	err = c.Set("test.credentials", `${CREDENTIALS | {}}`)
	assert.NoErr(t, err)

	// with env values
	t.Run("with env values", func(t *testing.T) {
		testutil.MockEnvValues(map[string]string{
			"CREDENTIALS": `{"username":"admin"}`,
		}, func() {
			st := Configuration{}
			err = c.Decode(&st)
			assert.NoErr(t, err)
			dump.Println(st)
			assert.Eq(t, `{"username":"admin"}`, st.Details.Credentials)
		})
	})
}

// https://github.com/gookit/config/issues/178
func TestIssues_178(t *testing.T) {
	type ConferenceConfigure struct {
		AuthServerEnable bool `mapstructure:"authServerEnable" default:"true"`
	}

	var ENVS = map[string]string{
		"CONF_AUTH_SERVER_ENABLE": "authServerEnable",
	}

	config.WithOptions(config.ParseEnv, config.ParseTime, config.ParseDefault)
	config.LoadOSEnvs(ENVS)

	cfg := &ConferenceConfigure{}
	err := config.Decode(cfg)
	assert.NoErr(t, err)
	dump.Println(cfg)
}

// https://github.com/gookit/config/issues/192
func TestIssues_192(t *testing.T) {
	s := `{
	"key": 23707729876828933003792990320594511132013137629744363463325945636682800546201191581706241551352734654762086038344743940857801503840360878427584255703013924373301145683882034301334533678253123777083489887967659929148298684008991665609773532863485728577470710590688325197694460521376123072613857785739366064688074459399408762960887169067851291291611970194076234580897060365318108861340336375060983779163595033605894055218557648363640361256922411962394084268360413547861005069585285713253756167043430574673046032573767256949834558011358364591391964981157578571244295339467926678988648036459748021538498339258036608168313
}`

	jsd := config.NewDriver("json", func(blob []byte, v any) (err error) {
		jnd := json.NewDecoder(bytes.NewReader(blob))
		jnd.UseNumber()
		return jnd.Decode(v)
	}, config.JSONEncoder)

	cfg := config.NewEmpty("test1").WithDriver(jsd)
	err := cfg.LoadStrings(config.JSON, s)
	assert.NoErr(t, err)
	dump.P(cfg.Data())

	// to big.Int
	bi := new(big.Int)
	_, ok := bi.SetString(cfg.String("key"), 10)
	assert.True(t, ok)
}

// https://github.com/gookit/config/issues/193 Support Environment Variable Overrides
func TestIssues_193(t *testing.T) {
	c := config.NewGeneric("test", config.WithTagName("config"))

	err := c.LoadStrings(config.JSON, `{"name": "default"}`)
	assert.NoErr(t, err)
	assert.Eq(t, "default", c.String("name"))

	err = c.LoadStrings(config.JSON, `{"datasource": {"username": "name in dev"}}`)
	assert.NoErr(t, err)
	assert.Eq(t, "name in dev", c.String("datasource.username"))

	testutil.MockEnvValue("DATASOURCE_USERNAME", "name in prod", func(val string) {
		c.LoadOSEnvs(map[string]string{
			"DATASOURCE_USERNAME": "datasource.username",
		})
	})

	assert.Eq(t, "name in prod", c.String("datasource.username"))
}

// https://github.com/gookit/config/issues/194
func TestIssues_194(t *testing.T) {
	cl := config.New("test", config.ParseDefault)

	type TestConfig struct {
		Nested struct {
			SimpleValue string
			WithDefault string `default:"default-value"`
		} `default:""` // <-- add this line
	}

	cfg := TestConfig{}
	err := cl.BindStruct("", &cfg)
	assert.NoErr(t, err)
	dump.P(cfg)
	assert.Eq(t, "default-value", cfg.Nested.WithDefault)
}

// https://github.com/gookit/config/issues/195
func TestIssues_195(t *testing.T) {
	type TestConfig struct {
		WithDefault time.Duration `default:"10s"`
	}

	defer config.Reset()
	config.WithOptions(
		config.ParseTime,
		config.ParseDefault,
	)
	config.AddDriver(yaml.Driver)
	err := config.LoadStrings(config.Yaml, `
name: inhere
`)
	assert.NoErr(t, err)

	cfg := TestConfig{}
	assert.NoErr(t, config.Decode(&cfg))
	dump.P(cfg)
}
