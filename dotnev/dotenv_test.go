package dotnev

import (
	"fmt"
	"os"
	"runtime"
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

	ClearLoaded()
}

func TestLoadExists(t *testing.T) {
	assert.Equal(t, "", os.Getenv("DONT_ENV_TEST"))

	err := LoadExists("./testdata", "not-exist", ".env")

	assert.NoError(t, err)
	assert.Equal(t, "blog", os.Getenv("DONT_ENV_TEST"))
	assert.Equal(t, "blog", Get("DONT_ENV_TEST"))
	ClearLoaded()
}

func TestLoadFromMap(t *testing.T) {
	assert.Equal(t, "", os.Getenv("DONT_ENV_TEST"))

	err := LoadFromMap(map[string]string{
		"DONT_ENV_TEST":  "blog",
		"dont_env_test1": "val1",
		"dont_env_test2": "23",
	})

	assert.NoError(t, err)

	envStr := fmt.Sprint(os.Environ())
	assert.Contains(t, envStr, "DONT_ENV_TEST=blog")
	assert.Contains(t, envStr, "DONT_ENV_TEST1=val1")

	assert.Equal(t, "blog", Get("DONT_ENV_TEST"))
	assert.Equal(t, "blog", os.Getenv("DONT_ENV_TEST"))
	assert.Equal(t, "val1", Get("DONT_ENV_TEST1"))
	assert.Equal(t, 23, Int("DONT_ENV_TEST2"))

	// on windows, os.Getenv() not case sensitive
	if runtime.GOOS == "windows" {
		assert.Equal(t, "val1", Get("dont_env_test1"))
		assert.Equal(t, 23, Int("dont_env_test2"))
	} else {
		assert.Equal(t, "val1", Get("dont_env_test1"))
		assert.Equal(t, 23, Int("dont_env_test2"))
	}

	assert.Equal(t, 20, Int("dont_env_test1", 20))
	assert.Equal(t, 20, Int("dont_env_not_exist", 20))

	// check cache
	assert.Contains(t, LoadedData(), "DONT_ENV_TEST2")

	// clear
	ClearLoaded()
	assert.Equal(t, "", os.Getenv("DONT_ENV_TEST"))
	assert.Equal(t, "", Get("DONT_ENV_TEST1"))

	err = LoadFromMap(map[string]string{
		"": "val",
	})
	assert.Error(t, err)
}

func TestDontUpperEnvKey(t *testing.T) {
	assert.Equal(t, "", os.Getenv("DONT_ENV_TEST"))

	DontUpperEnvKey()

	err := LoadFromMap(map[string]string{
		"dont_env_test": "val",
	})

	assert.Contains(t, fmt.Sprint(os.Environ()), "dont_env_test=val")
	assert.NoError(t, err)
	assert.Equal(t, "val", Get("dont_env_test"))

	// on windows, os.Getenv() not case sensitive
	if runtime.GOOS == "windows" {
		assert.Equal(t, "val", Get("DONT_ENV_TEST"))
	} else {
		assert.Equal(t, "", Get("DONT_ENV_TEST"))
	}

	UpperEnvKey = true // revert
	ClearLoaded()
}
