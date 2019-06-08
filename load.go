package config

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/imdario/mergo"
)

// LoadFiles load one or multi files
func LoadFiles(sourceFiles ...string) error { return dc.LoadFiles(sourceFiles...) }

// LoadFiles load and parse config files
func (c *Config) LoadFiles(sourceFiles ...string) (err error) {
	for _, file := range sourceFiles {
		if err = c.loadFile(file, false); err != nil {
			return
		}
	}
	return
}

// LoadExists load one or multi files, will ignore not exist
func LoadExists(sourceFiles ...string) error { return dc.LoadExists(sourceFiles...) }

// LoadExists load and parse config files, but will ignore not exists file.
func (c *Config) LoadExists(sourceFiles ...string) (err error) {
	for _, file := range sourceFiles {
		if err = c.loadFile(file, true); err != nil {
			return
		}
	}
	return
}

// LoadRemote load config data from remote URL.
func LoadRemote(format, url string) error { return dc.LoadRemote(format, url) }

// LoadRemote load config data from remote URL.
// Usage:
// 	c.LoadRemote(config.JSON, "http://abc.com/api-config.json")
func (c *Config) LoadRemote(format, url string) (err error) {
	// create http client
	client := http.Client{Timeout: 300 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("fetch remote resource error, reply status code is not equals to 200")
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

// LoadOSEnv load data from OS ENV
func LoadOSEnv(keys []string) { dc.LoadOSEnv(keys) }

// LoadOSEnv load data from os ENV
func (c *Config) LoadOSEnv(keys []string) {
	for _, key := range keys {
		val := os.Getenv(strings.ToUpper(key))
		_ = c.Set(key, val)
	}
}

// support bound types for CLI flags vars
var validTypes = map[string]int{
	"int":  1,
	"uint": 1,
	"bool": 1,
	// string is default
	"string": 1,
}

// LoadFlags load data from cli flags
func LoadFlags(keys []string) error { return dc.LoadFlags(keys) }

// LoadFlags parse command line arguments, based on provide keys.
// Usage:
// 	// debug flag is bool type
// 	c.LoadFlags([]string{"env", "debug:bool"})
func (c *Config) LoadFlags(keys []string) (err error) {
	hash := map[string]interface{}{}

	// bind vars
	for _, key := range keys {
		key, typ := parseVarNameAndType(key)

		switch typ {
		case "int":
			ptr := new(int)
			flag.IntVar(ptr, key, c.Int(key), "")
			hash[key] = ptr
		case "uint":
			ptr := new(uint)
			flag.UintVar(ptr, key, c.Uint(key), "")
			hash[key] = ptr
		case "bool":
			ptr := new(bool)
			flag.BoolVar(ptr, key, c.Bool(key), "")
			hash[key] = ptr
		default: // as string
			ptr := new(string)
			flag.StringVar(ptr, key, c.String(key), "")
			hash[key] = ptr
		}
	}

	// parse and collect
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		name := f.Name
		// name := strings.Replace(f.Name, "-", ".", -1)
		// only get name in the keys.
		if _, ok := hash[name]; !ok {
			return
		}

		// ignore error
		_ = c.Set(name, f.Value.String())
	})
	return
}

func parseVarNameAndType(key string) (string, string) {
	typ := "string"
	key = strings.Trim(key, "-")

	// can set var type: int, uint, bool
	if strings.IndexByte(key, ':') > 0 {
		list := strings.SplitN(key, ":", 2)
		key, typ = list[0], list[1]

		if _, ok := validTypes[typ]; !ok {
			typ = "string"
		}
	}
	return key, typ
}

// LoadData load one or multi data
func LoadData(dataSource ...interface{}) error { return dc.LoadData(dataSource...) }

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

// LoadSources load one or multi byte data
func LoadSources(format string, src []byte, more ...[]byte) error {
	return dc.LoadSources(format, src, more...)
}

// LoadSources load data from byte content.
// Usage:
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

// LoadStrings load one or multi string
func LoadStrings(format string, str string, more ...string) error {
	return dc.LoadStrings(format, str, more...)
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
