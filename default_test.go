package config

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultLoad(t *testing.T) {
	st := assert.New(t)

	ClearAll()
	err := LoadFiles("testdata/json_base.json", "testdata/json_other.json")
	st.Nil(err)

	ClearAll()
	err = LoadExists("testdata/json_base.json", "not-exist.json")
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

func TestSetDecoderEncoder(t *testing.T) {
	at := assert.New(t)

	c := Default()
	c.ClearAll()

	at.True(c.HasDecoder(Json))
	at.True(c.HasEncoder(Json))

	c.DelDriver(Json)

	at.False(c.HasDecoder(Json))
	at.False(c.HasEncoder(Json))

	SetDecoder(Json, JsonDecoder)
	SetEncoder(Json, JsonEncoder)

	at.True(c.HasDecoder(Json))
	at.True(c.HasEncoder(Json))
}

func TestDefault(t *testing.T) {
	at := assert.New(t)

	ClearAll()
	dc.initialized = false
	WithOptions(ParseEnv)

	at.True(GetOptions().ParseEnv)

	LoadStrings(Json, `{"name": "inhere"}`)

	buf := &bytes.Buffer{}
	_, err := WriteTo(buf)
	at.Nil(err)
}
