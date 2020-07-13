package config

import (
	"bytes"
	"testing"

	"github.com/gookit/goutil/dump"
	"github.com/stretchr/testify/assert"
)

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

func TestConfig_Structure(t *testing.T) {
	is := assert.New(t)

	cfg := Default()
	cfg.ClearAll()

	err := cfg.LoadStrings(JSON, `{
"age": 28,
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)

	is.Nil(err)

	user := &struct {
		Age    int // always float64 from JSON
		Name   string
		Sports []string
	}{}
	// map all data
	err = MapStruct("", user)
	is.Nil(err)

	is.Equal(28, user.Age)
	is.Equal("inhere", user.Name)
	is.Equal("pingPong", user.Sports[0])

	// - auto convert string to int
	// age use string in JSON
	cfg1 := New("test")
	err = cfg1.LoadStrings(JSON, `{
"age": "26",
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)

	is.Nil(err)

	user1 := &struct {
		Age    int // always float64 from JSON
		Name   string
		Sports []string
	}{}
	err = cfg1.MapStruct("", user1)
	is.Nil(err)

	dump.P(*user1)

	// map some data
	err = cfg.LoadStrings(JSON, `{
"sec": {
	"key": "val",
	"age": 120,
	"tags": [12, 34]
}
}`)
	is.Nil(err)

	some := struct {
		Age  int
		Kye  string
		Tags []int
	}{}
	err = BindStruct("sec", &some)
	is.Nil(err)
	is.Equal(120, some.Age)
	is.Equal(12, some.Tags[0])
	cfg.ClearAll()

	// custom data
	cfg = New("test")
	err = cfg.LoadData(map[string]interface{}{
		"key":  "val",
		"age":  120,
		"tags": []int{12, 34},
	})
	is.NoError(err)

	s1 := struct {
		Age  int
		Kye  string
		Tags []int
	}{}
	err = cfg.BindStruct("", &s1)
	is.Nil(err)
	is.Equal(120, s1.Age)
	is.Equal(12, s1.Tags[0])

	// key not exist
	err = cfg.BindStruct("not-exist", &s1)
	is.Error(err)
	is.Equal("this key does not exist in the config", err.Error())

	cfg.ClearAll()
}
