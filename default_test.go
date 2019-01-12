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

	at.True(c.HasDecoder(JSON))
	at.True(c.HasEncoder(JSON))

	c.DelDriver(JSON)

	at.False(c.HasDecoder(JSON))
	at.False(c.HasEncoder(JSON))

	SetDecoder(JSON, JSONDecoder)
	SetEncoder(JSON, JSONEncoder)

	at.True(c.HasDecoder(JSON))
	at.True(c.HasEncoder(JSON))
}

func TestDefault(t *testing.T) {
	at := assert.New(t)

	ClearAll()
	WithOptions(ParseEnv)

	at.True(GetOptions().ParseEnv)

	_ = LoadStrings(JSON, `{"name": "inhere"}`)

	buf := &bytes.Buffer{}
	_, err := WriteTo(buf)
	at.Nil(err)
}
