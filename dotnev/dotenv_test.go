package dotnev

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	err := Load("./testdata", "not-exist", ".env")
	assert.Error(t, err)

	assert.Equal(t, "", os.Getenv("APP_NAME"))

	err = Load("./testdata")
	assert.NoError(t, err)
	assert.Equal(t, "blog", os.Getenv("APP_NAME"))
	assert.Equal(t, "blog", Get("APP_NAME"))
	_ = os.Unsetenv("APP_NAME")

	err = Load("./testdata", "error.ini")
	assert.Error(t, err)

	err = Load("./testdata", "invalid_key.ini")
	assert.Error(t, err)

	assert.Equal(t, "def-val", Get("NOT-EXIST", "def-val"))
}

func TestLoadExists(t *testing.T) {
	assert.Equal(t, "", os.Getenv("APP_NAME"))

	err := LoadExists("./testdata", "not-exist", ".env")

	assert.NoError(t, err)
	assert.Equal(t, "blog", os.Getenv("APP_NAME"))
	assert.Equal(t, "blog", Get("APP_NAME"))
	_ = os.Unsetenv("APP_NAME")
}

func TestLoadFromMap(t *testing.T) {
	assert.Equal(t, "", os.Getenv("APP_NAME"))

	err := LoadFromMap(map[string]string{
		"APP_NAME": "blog",
	})

	assert.NoError(t, err)
	assert.Equal(t, "blog", os.Getenv("APP_NAME"))
	_ = os.Unsetenv("APP_NAME")

	err = LoadFromMap(map[string]string{
		"": "val",
	})
	assert.Error(t, err)
}
