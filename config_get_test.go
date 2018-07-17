package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	st := assert.New(t)

	ClearAll()
	err := LoadStrings(Json, jsonStr)
	st.Nil(err)

	// get bool
	bv, ok := Get("debug")
	st.True(ok)
	st.Equal(true, bv)

	// get string
	val, ok := Get("name")
	st.True(ok)
	st.Equal("app", val)
}
