package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ast := assert.New(t)

	err := LoadFiles("testdata/json_base.json")
	if err != nil {
		t.Error(err)
	}

	val, ok := Get("name")
	ast.True(ok)
	ast.Equal("app", val)
}
