// Package dotnev provide load .env data to os ENV
package dotnev

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/ini/v2/parser"
)

// DefaultName default file name
var DefaultName = ".env"

// OnlyLoadExists load on file exists
var OnlyLoadExists bool

// Load parse .env file data to os ENV.
// Usage:
// 	dotenv.Load("./", ".env")
func Load(dir string, filenames ...string) (err error) {
	if len(filenames) == 0 {
		filenames = []string{DefaultName}
	}

	for _, filename := range filenames {
		file := filepath.Join(dir, filename)
		if err = loadFile(file); err != nil {
			break
		}
	}
	return
}

// LoadExists load on file exists
func LoadExists(dir string, filenames ...string) error {
	OnlyLoadExists = true

	return Load(dir, filenames...)
}

// LoadFromMap load data from given string map
func LoadFromMap(kv map[string]string) (err error) {
	for key, val := range kv {
		key = strings.ToUpper(key)
		err = os.Setenv(key, val)
		if err != nil {
			break
		}
	}
	return
}

// Get get os ENV value by name
func Get(name string, defVal ...string) (val string) {
	name = strings.ToUpper(name)
	if val = os.Getenv(name); val != "" {
		return
	}

	if len(defVal) > 0 {
		val = defVal[0]
	}
	return
}

// load and parse .env file data to os ENV
func loadFile(file string) (err error) {
	// open file
	fd, err := os.Open(file)
	if err != nil {
		// skip not exist file
		if os.IsNotExist(err) && OnlyLoadExists {
			return nil
		}
		return err
	}
	defer fd.Close()

	// parse file content
	s := bufio.NewScanner(fd)
	p := parser.NewSimpled(parser.NoDefSection)

	if _, err = p.ParseFrom(s); err != nil {
		return
	}

	// set data to os ENV
	if mp, ok := p.SimpleData()[p.DefSection]; ok {
		err = LoadFromMap(mp)
	}
	return
}
