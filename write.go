package config

import (
	"errors"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/maputil"
)

// some common errors
var (
	ErrReadonly   = errors.New("the config instance in 'readonly' mode")
	ErrKeyIsEmpty = errors.New("the config key is cannot be empty")
)

// SetData for override the Config.Data
func SetData(data map[string]interface{}) {
	dc.SetData(data)
}

// SetData for override the Config.Data
func (c *Config) SetData(data map[string]interface{}) {
	c.lock.Lock()
	c.data = data
	c.lock.Unlock()

	c.fireHook(OnSetData)
}

// Set val by key
func Set(key string, val interface{}, setByPath ...bool) error {
	return dc.Set(key, val, setByPath...)
}

// Set a value by key string.
func (c *Config) Set(key string, val interface{}, setByPath ...bool) (err error) {
	if c.opts.Readonly {
		return ErrReadonly
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	sep := c.opts.Delimiter
	if key = formatKey(key, string(sep)); key == "" {
		return ErrKeyIsEmpty
	}

	defer c.fireHook(OnSetValue)
	if strings.IndexByte(key, sep) == -1 {
		c.data[key] = val
		return
	}

	// disable set by path.
	if len(setByPath) > 0 && !setByPath[0] {
		c.data[key] = val
		return
	}

	// set by path
	keys := strings.Split(key, string(sep))
	return maputil.SetByKeys(&c.data, keys, val)
}

/**
more setter: SetIntArr, SetIntMap, SetString, SetStringArr, SetStringMap
*/

// build new value by key paths
// "site.info" -> map[string]map[string]val
func buildValueByPath(paths []string, val interface{}) (newItem map[string]interface{}) {
	if len(paths) == 1 {
		return map[string]interface{}{paths[0]: val}
	}

	arrutil.Reverse(paths)

	// multi nodes
	for _, p := range paths {
		if newItem == nil {
			newItem = map[string]interface{}{p: val}
		} else {
			newItem = map[string]interface{}{p: newItem}
		}
	}
	return
}

// build new value by key paths, only for yaml.v2
// "site.info" -> map[interface{}]map[string]val
func buildValueByPath1(paths []string, val interface{}) (newItem map[interface{}]interface{}) {
	if len(paths) == 1 {
		return map[interface{}]interface{}{paths[0]: val}
	}

	arrutil.Reverse(paths)

	// multi nodes
	for _, p := range paths {
		if newItem == nil {
			newItem = map[interface{}]interface{}{p: val}
		} else {
			newItem = map[interface{}]interface{}{p: newItem}
		}
	}
	return
}
