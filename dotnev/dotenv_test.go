package dotnev

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	err := Load("./testdata", "not-exist", ".env")
	assert.Error(t, err)

	assert.Equal(t, "", os.Getenv("DONT_ENV_TEST"))

	err = Load("./testdata")
	assert.NoError(t, err)
	assert.Equal(t, "blog", os.Getenv("DONT_ENV_TEST"))
	assert.Equal(t, "blog", Get("DONT_ENV_TEST"))
	_ = os.Unsetenv("DONT_ENV_TEST") // clear

	err = Load("./testdata", "error.ini")
	assert.Error(t, err)

	err = Load("./testdata", "invalid_key.ini")
	assert.Error(t, err)

	assert.Equal(t, "def-val", Get("NOT-EXIST", "def-val"))
}

func TestLoadExists(t *testing.T) {
	assert.Equal(t, "", os.Getenv("DONT_ENV_TEST"))

	err := LoadExists("./testdata", "not-exist", ".env")

	assert.NoError(t, err)
	assert.Equal(t, "blog", os.Getenv("DONT_ENV_TEST"))
	assert.Equal(t, "blog", Get("DONT_ENV_TEST"))
	_ = os.Unsetenv("DONT_ENV_TEST") // clear
}

func TestLoadFromMap(t *testing.T) {
	assert.Equal(t, "", os.Getenv("DONT_ENV_TEST"))

	err := LoadFromMap(map[string]string{
		"DONT_ENV_TEST": "blog",
		"dont_env_test1": "val1",
	})

	assert.NoError(t, err)
	assert.Equal(t, "blog", Get("DONT_ENV_TEST"))
	assert.Equal(t, "val1", os.Getenv("DONT_ENV_TEST1"))
	assert.Equal(t, "val1", Get("dont_env_test1"))
	// clear
	_ = os.Unsetenv("DONT_ENV_TEST")
	_ = os.Unsetenv("dont_env_test1")

	err = LoadFromMap(map[string]string{
		"": "val",
	})
	assert.Error(t, err)
}

func TestDontUpperEnvKey(t *testing.T) {
	DontUpperEnvKey()

	err := LoadFromMap(map[string]string{
		"dont_env_test": "blog",
	})
	assert.NoError(t, err)
	assert.Equal(t, "blog", Get("DONT_ENV_TEST"))

	UpperEnvKey = true // revert
}