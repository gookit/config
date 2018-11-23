package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/imdario/mergo"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LoadFiles load and parse config files
func (c *Config) LoadFiles(sourceFiles ...string) (err error) {
	for _, file := range sourceFiles {
		err = c.loadFile(file, false)
		if err != nil {
			return
		}
	}

	c.initialized = true
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

	c.initialized = true
	return
}

// load config file
func (c *Config) loadFile(file string, loadExist bool) (err error) {
	// open file
	fd, err := os.Open(file)
	if err != nil {
		// skip not exist file
		if os.IsNotExist(err) && loadExist {
			return nil
		}

		return err
	}
	defer fd.Close()

	// read file content
	bts, err := ioutil.ReadAll(fd)
	if err == nil {
		// get format for file ext
		format := strings.Trim(filepath.Ext(file), ".")

		// parse file content
		if err = c.parseSourceCode(format, bts); err != nil {
			return
		}

		c.loadedFiles = append(c.loadedFiles, file)
	}

	return
}

// LoadRemote load config data from remote URL.
// Usage:
// 	c.LoadRemote(config.JSON, "http://abc.com/api-config.json")
func (c *Config) LoadRemote(format, url string) (err error) {
	// create http client
	client := http.Client{Timeout: 900 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("remote resource is not exist, reply status code is not equals to 200")
	}

	// read response content
	bts, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		// parse file content
		if err = c.parseSourceCode(format, bts); err != nil {
			return
		}

		c.loadedFiles = append(c.loadedFiles, url)
	}

	return
}

// LoadFlags parse command line arguments, based on provide keys.
// Usage:
// 	c.LoadFlags([]string{"env", "debug"})
func (c *Config) LoadFlags(keys []string) (err error) {
	hash := map[string]*string{}
	for _, key := range keys {
		hash[key] = new(string)
		defVal, _ := c.String(key)
		flag.StringVar(hash[key], key, defVal, "")
	}

	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		name := f.Name
		// name := strings.Replace(f.Name, "-", ".", -1)

		// only get name in the keys.
		if _, ok := hash[name]; !ok {
			return
		}

		err = c.Set(name, f.Value.String())
		if err != nil {
			return
		}
	})

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

	c.initialized = true
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

	c.initialized = true
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

	c.initialized = true
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
	case JSON:
		decoder, ok = c.decoders[JSON]
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
