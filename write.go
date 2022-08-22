package config

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/imdario/mergo"
)

var (
	errReadonly   = errors.New("the config instance in 'readonly' mode")
	errKeyIsEmpty = errors.New("the config key is cannot be empty")
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
		return errReadonly
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	sep := c.opts.Delimiter
	if key = formatKey(key, string(sep)); key == "" {
		return errKeyIsEmpty
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

	keys := strings.Split(key, string(sep))
	topK := keys[0]
	paths := keys[1:]

	var ok bool
	var item interface{}

	// find top item data based on top key
	if item, ok = c.data[topK]; !ok {
		// not found, is new add
		c.data[topK] = buildValueByPath(paths, val)
		return
	}

	switch typeData := item.(type) {
	case map[interface{}]interface{}: // from yaml.v2
		dstItem := make(map[interface{}]interface{}, len(typeData))
		for k, v := range typeData {
			dstItem[strutil.QuietString(k)] = v
		}

		// create a new item for the topK
		newItem := buildValueByPath1(paths, val)
		// merge new item to old item
		if err = mergo.Merge(&dstItem, newItem, mergo.WithOverride); err != nil {
			return
		}

		c.data[topK] = dstItem
	case map[string]interface{}: // from json,toml,yaml.v3
		// enhanced: 'top.sub'=string, can set 'top.sub' to other type. eg: 'top.sub'=map
		if err = maputil.SetByKeys(&typeData, paths, val); err != nil {
			return
		}

		c.data[topK] = typeData
	case []interface{}: // is list array
		index, err := strconv.Atoi(keys[1])
		if len(keys) == 2 && err == nil {
			if index <= len(typeData) {
				typeData[index] = val
			}

			c.data[topK] = typeData
		} else {
			err = errors.New("max allow 1 level for setting array value, current key: " + key)
			return err
		}
	default:
		// as a top key
		c.data[key] = val
		// err = errors.New("not supported value type, cannot setting value for the key: " + key)
	}
	return
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
