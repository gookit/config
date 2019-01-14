// Package dotnev provide load .env data to os ENV
package dotnev

import (
	"bufio"
	"github.com/gookit/ini/parser"
	"os"
	"path/filepath"
	"strings"
)

// LoadExists Load data on file exists
var LoadExists bool

// Load parse .env file data to os ENV.
// Usage:
//	dotenv.Load("./", ".env")
func Load(dir string, filenames ...string) (err error) {
	for _, filename := range filenames {
		file := filepath.Join(dir, filename)
		if err = loadFile(file); err != nil {
			break
		}
	}
	return
}

// load and parse .env file data to os ENV
func loadFile(file string) (err error) {
	// open file
	fd, err := os.Open(file)
	if err != nil {
		// skip not exist file
		if os.IsNotExist(err) && LoadExists {
			return nil
		}
		return err
	}
	defer fd.Close()

	// parse file content
	s := bufio.NewScanner(fd)
	p := parser.SimpleParser(parser.NoDefSection)

	if _, err = p.ParseFrom(s); err != nil {
		return
	}

	// set data to os ENV
	if mp, ok := p.SimpleData()[p.DefSection]; ok {
		for key, val := range mp {
			err = os.Setenv(strings.ToUpper(key), val)
			if err != nil {
				break
			}
		}
	}
	return
}
