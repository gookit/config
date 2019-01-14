package dotnev

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	err := Load("./testdata", "not-exist", ".env")
	assert.Error(t, err)

	assert.Equal(t, "", os.Getenv("APP_NAME"))

	err = Load("./testdata")
	assert.NoError(t, err)
	assert.Equal(t, "blog", os.Getenv("APP_NAME"))
	_= os.Unsetenv("APP_NAME")

	err = Load("./testdata", "error.ini")
	assert.Error(t, err)

	err = Load("./testdata", "invalid_key.ini")
	assert.Error(t, err)
}

func TestLoadExists(t *testing.T) {
	assert.Equal(t, "", os.Getenv("APP_NAME"))

	err := LoadExists("./testdata", "not-exist", ".env")
	assert.NoError(t, err)
	assert.Equal(t, "blog", os.Getenv("APP_NAME"))
	_= os.Unsetenv("APP_NAME")
}
