package config

import (
	"bytes"
	"testing"

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
	st := assert.New(t)

	cfg := Default()
	cfg.ClearAll()

	err := cfg.LoadStrings(JSON, `{
"age": 28,
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)

	st.Nil(err)

	user := &struct {
		Age    int // always float64 from JSON
		Name   string
		Sports []string
	}{}
	// map all data
	err = MapStruct("", user)
	st.Nil(err)

	st.Equal(28, user.Age)
	st.Equal("inhere", user.Name)
	st.Equal("pingPong", user.Sports[0])

	// map some data
	err = cfg.LoadStrings(JSON, `{
"sec": {
	"key": "val",
	"age": 120,
	"tags": [12, 34]
}
}`)
	st.Nil(err)

	some := struct {
		Age  int
		Kye  string
		Tags []int
	}{}
	err = BindStruct("sec", &some)
	st.Nil(err)
	st.Equal(120, some.Age)
	st.Equal(12, some.Tags[0])
	cfg.ClearAll()

	// custom data
	cfg = New("test")
	err = cfg.LoadData(map[string]interface{}{
		"key":  "val",
		"age":  120,
		"tags": []int{12, 34},
	})
	st.NoError(err)

	s1 := struct {
		Age  int
		Kye  string
		Tags []int
	}{}
	err = cfg.BindStruct("", &s1)
	st.Nil(err)
	st.Equal(120, s1.Age)
	st.Equal(12, s1.Tags[0])

	// key not exist
	err = cfg.BindStruct("not-exist", &s1)
	st.Error(err)
	st.Equal("this key does not exist in the config", err.Error())

	cfg.ClearAll()
}
