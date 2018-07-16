package config

import (
	"strings"
	"errors"
	"github.com/imdario/mergo"
)

// Set a value by key string.
func (c *Config) Set(key string, val interface{}, setByPath ...bool) (err error) {
	// if is readonly
	if c.opts.Readonly {
		err = errors.New("the config instance in 'readonly' mode")
		return
	}

	key = strings.Trim(strings.TrimSpace(key), ".")
	if key == "" {
		err = errors.New("the config key is cannot be empty")
		return
	}

	// is top key
	if !strings.Contains(key, ".") {
		c.data[key] = val
		return
	}

	// disable set by path.
	if len(setByPath) > 0 && !setByPath[0] {
		c.data[key] = val
		return
	}

	keys := strings.Split(key, ".")
	topK := keys[0]

	var ok bool
	var item interface{}

	// find top item data based on top key
	if item, ok = c.data[topK]; !ok {
		// not found, is new add
		c.data[topK] = buildValueByPath(keys[1:], val)
		return
	}

	// create a new item for the topK
	newItem := buildValueByPath(keys[1:], val)

	// merge new item to old item
	err = mergo.Map(&item, newItem, mergo.WithOverride)
	if err != nil {
		return
	}

	// resetting
	c.data[topK] = item
	return
}

// build new value by key paths
// "site.info" -> map[string]map[string]val
func buildValueByPath(paths []string, val interface{}) interface{} {
	if len(paths) == 1 {
		if paths[0] == "0" {
			return []interface{}{val}
		} else {
			return map[string]interface{}{paths[0]: val}
		}
	}

	sliceReverse(paths)

	// multi nodes
	for _, p := range paths {
		if p == "0" {
			val = []interface{}{val}
		} else {
			val = map[string]interface{}{p: val}
		}
	}

	return val
}

// reverse a slice. (slice 是引用，所以可以直接改变)
func sliceReverse(ss []string) {
	ln := len(ss)

	for i := 0; i < int(ln/2); i++ {
		li := ln - i - 1
		// fmt.Println(i, "<=>", li)
		ss[i], ss[li] = ss[li], ss[i]
	}
}
