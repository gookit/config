package config

import (
	"strings"
	"fmt"
	"errors"
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
		err = errors.New(fmt.Sprintf("the top key '%s' is not exist", topK))
		return
	}

	// find child
	for _, k := range keys[1:] {
		switch item.(type) {
		case int, int8, int64, string, bool:
			err = errors.New("the top item is scalar type, cannot set sub-value by path: " + key)
			return
		case map[string]interface{}:
			item, ok = item.(map[string]interface{})[k]
			if !ok {
				return
			}
		}
	}

	return
}

// build new value by key paths
// "site.info" -> map[string]map[string]val
func buildValueByPath(paths []string, val interface{}) interface{} {
	ln := len(paths)
	sliceReverse(paths)

	if ln == 1 {
		if paths[0] == "0" {
			return []interface{}{val}
		} else {
			return map[string]interface{}{paths[0]: val}
		}
	}

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
