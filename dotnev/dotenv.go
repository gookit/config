// Package dotnev provide load .env data to os ENV
//
// Deprecated: please use github.com/gookit/ini/v2/dotenv
package dotnev

import (
	"github.com/gookit/ini/v2/dotenv"
)

// LoadedData get all loaded data by dontenv
// Deprecated: please use github.com/gookit/ini/v2/dotenv
func LoadedData() map[string]string {
	return dotenv.LoadedData()
}

// ClearLoaded clear the previously set ENV value
// Deprecated: please use github.com/gookit/ini/v2/dotenv
func ClearLoaded() {
	dotenv.ClearLoaded()
}

// DontUpperEnvKey don't change key to upper on set ENV
// Deprecated: please use github.com/gookit/ini/v2/dotenv
func DontUpperEnvKey() {
	dotenv.DontUpperEnvKey()
}

// Load parse .env file data to os ENV.
//
// Usage:
//
//	dotenv.Load("./", ".env")
//
// Deprecated: please use github.com/gookit/ini/v2/dotenv
func Load(dir string, filenames ...string) (err error) {
	return dotenv.Load(dir, filenames...)
}

// LoadExists only load on file exists
// Deprecated: please use github.com/gookit/ini/v2/dotenv
func LoadExists(dir string, filenames ...string) error {
	return dotenv.LoadExists(dir, filenames...)
}

// LoadFromMap load data from given string map
// Deprecated: please use github.com/gookit/ini/v2/dotenv
func LoadFromMap(kv map[string]string) (err error) {
	return dotenv.LoadFromMap(kv)
}

// Get get os ENV value by name
// Deprecated: please use github.com/gookit/ini/v2/dotenv
func Get(name string, defVal ...string) (val string) {
	return dotenv.Get(name, defVal...)
}

// Bool get a bool value by key
// Deprecated: please use github.com/gookit/ini/v2/dotenv
func Bool(name string, defVal ...bool) (val bool) {
	return dotenv.Bool(name, defVal...)
}

// Int get a int value by key
// Deprecated: please use github.com/gookit/ini/v2/dotenv
func Int(name string, defVal ...int) (val int) {
	return dotenv.Int(name, defVal...)
}
