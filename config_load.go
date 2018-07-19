package config

import (
	"os"
	"strings"
	"path/filepath"
	"io/ioutil"
	"github.com/imdario/mergo"
	"errors"
)

// LoadFiles load and parse config files
func (c *Config) LoadFiles(sourceFiles ...string) (err error) {
	for _, file := range sourceFiles {
		err = c.loadFile(file, false)
		if err != nil {
			return
		}
	}

	return
}

// LoadExists load and parse config files, but will ignore not exists file.
func (c *Config) LoadExists(sourceFiles ...string) (err error) {
	for _, file := range sourceFiles {
		err = c.loadFile(file, true)
		if err != nil {
			return
		}
	}

	return
}

// load config file
func (c *Config) loadFile(file string, loadExist bool) (err error) {
	if _, err = os.Stat(file); err != nil {
		// skip not exist file
		if os.IsNotExist(err) && loadExist {
			return nil
		}

		return
	}

	// open file
	fd, err := os.Open(file)
	if err != nil {
		return
	}
	defer fd.Close()

	// read file content
	content, err := ioutil.ReadAll(fd)
	if err == nil {
		// get format for file ext
		format := strings.Trim(filepath.Ext(file), ".")

		// parse file content
		if err = c.parseSourceCode(format, content); err != nil {
			return
		}

		c.loadedFiles = append(c.loadedFiles, file)
	}

	return
}

// LoadData load data from map OR struct
func (c *Config) LoadData(dataSources ...interface{}) (err error) {
	for _, ds := range dataSources {
		err = mergo.Merge(&c.data, ds, mergo.WithOverride)
		if err != nil {
			return
		}
	}

	return
}

// LoadSources load data from byte content.
// usage:
// 	config.LoadSources(config.Yml, []byte(`
// 	name: blog
// 	arr:
// 		key: val
// `))
func (c *Config) LoadSources(format string, src []byte, more ...[]byte) (err error) {
	err = c.parseSourceCode(format, src)
	if err != nil {
		return
	}

	for _, sc := range more {
		err = c.parseSourceCode(format, sc)
		if err != nil {
			return
		}
	}

	return
}

// LoadStrings load data from source string content.
func (c *Config) LoadStrings(format string, str string, more ...string) (err error) {
	err = c.parseSourceCode(format, []byte(str))
	if err != nil {
		return
	}

	for _, s := range more {
		err = c.parseSourceCode(format, []byte(s))
		if err != nil {
			return
		}
	}

	return
}

// parse config source code to Config.
func (c *Config) parseSourceCode(format string, blob []byte) (err error) {
	var ok bool
	var decoder Decoder

	switch format {
	case Hcl:
		decoder, ok = c.decoders[Hcl]
	case Ini:
		decoder, ok = c.decoders[Ini]
	case Json:
		decoder, ok = c.decoders[Json]
	case Yaml, Yml:
		decoder, ok = c.decoders[Yaml]
	case Toml:
		decoder, ok = c.decoders[Toml]
	}

	if !ok {
		return errors.New("no exists or no register decoder for the format: " + format)
	}

	data := make(map[string]interface{})

	// decode content to data
	if err = decoder(blob, &data); err != nil {
		return
	}

	// init config data
	if len(c.data) == 0 {
		c.data = data
	} else {
		// again ... will merge data
		// err = mergo.Map(&c.data, data, mergo.WithOverride)
		err = mergo.Merge(&c.data, data, mergo.WithOverride)
	}

	data = nil
	return
}
