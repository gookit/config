package config

import (
	"os"
	"strings"
	"path/filepath"
	"log"
	"io/ioutil"
	"github.com/imdario/mergo"
)

// LoadFiles load and parse config files
func (c *Config) LoadFiles(sourceFiles ...string) (err error) {
	for _, file := range sourceFiles {
		c.loadFile(file, false)
	}

	return
}

// LoadExists load and parse config files, but will ignore not exists file.
func (c *Config) LoadExists(sourceFiles ...string) (err error) {
	for _, file := range sourceFiles {
		c.loadFile(file, true)
	}

	return
}

// load config file
func (c *Config) loadFile(file string, onlyExist bool) (err error) {
	if _, err = os.Stat(file); err != nil {
		if os.IsNotExist(err) && onlyExist {
			return
		}

		panic(err)
	}

	// get format for file ext
	format := strings.Trim(filepath.Ext(file), ".")
	fd, err := os.Open(file)

	if err != nil {
		panic(err)
	}
	defer fd.Close()

	content, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(err)
	}

	// parse file content
	if c.parseSourceCode(format, content) != nil {
		return
	}

	c.loadedFiles = append(c.loadedFiles, file)

	return
}

// LoadData load data from map OR struct
func (c *Config) LoadData(dataSources ...interface{}) (err error) {
	for _, ds := range dataSources {
		err = mergo.Merge(&c.data, ds, mergo.WithOverride)
		if err != nil {
			panic(err)
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
func (c *Config) LoadSources(format string, sourceCodes ...[]byte) (err error) {
	for _, sc := range sourceCodes {
		err = c.parseSourceCode(format, sc)

		if err != nil {
			panic(err)
		}
	}

	return
}

// parse config source code to Config.
func (c *Config) parseSourceCode(format string, blob []byte) (err error) {
	var ok bool
	var decoder Decoder

	switch format {
	case Json:
		decoder, ok = c.decoders[Json]
	case Yaml, Yml:
		decoder, ok = c.decoders[Yaml]
	case Toml:
		decoder, ok = c.decoders[Toml]
	}

	if !ok {
		log.Fatalf("no exists or no register decoder for the format: %s", format)
	}

	data := make(map[string]interface{})

	// decode content to data
	if decoder(blob, &data) != nil {
		return
	}

	// init config data
	if len(c.data) == 0 {
		c.data = data

		// second... merge data
	} else {
		// err = mergo.Map(&c.data, data, mergo.WithOverride)
		err = mergo.Merge(&c.data, data, mergo.WithOverride)
	}

	return
}
