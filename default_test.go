package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestDefaultLoad(t *testing.T) {
	st := assert.New(t)

	ClearAll()
	err := LoadExists("testdata/json_base.json", "not-exist.json")
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
