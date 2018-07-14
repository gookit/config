package config

import (
	"strings"
	"strconv"
	"encoding/json"
	"fmt"
	"os"
)

// Get config value by key string, support get sub-value by key path(eg. 'map.key'),
// ok is true, find value from config
// ok is false, not found or error
func (c *Config) Get(key string, findByPath ...bool) (value interface{}, ok bool) {
	key = strings.Trim(strings.TrimSpace(key), ".")
	if key == "" {
		return
	}

	// check top key
	if value, ok = c.data[key]; ok {
		return
	}

	// disable find by path.
	if len(findByPath) > 0 && !findByPath[0] {
		return
	}

	// has sub key? eg. "lang.dir"
	if !strings.Contains(key, ".") {
		return
	}

	keys := strings.Split(key, ".")
	topK := keys[0]

	// 根据top key 找到顶级 item 数据
	var item interface{}
	if item, ok = c.data[topK]; !ok {
		return
	}

	// 查找子级
	for _, k := range keys[1:] {
		switch item.(type) {
		case map[interface{}]interface{}: // is map
			item, ok = item.(map[interface{}]interface{})[k]

			if !ok {
				return
			}
		case []interface{}: // is array
			i, err := strconv.Atoi(k)
			if err != nil {
				return
			}

			// 检查slice index是否存在
			arrItem := item.([]interface{})
			if len(arrItem) < i {
				return
			}

			item = arrItem[i]
		default: // error
			ok = false
			return
		}
	}

	return item, true
}

// GetString
func (c *Config) GetString(key string) (value string, ok bool) {
	val, ok := c.Get(key)
	if !ok {
		return
	}

	switch val.(type) {
	case bool, int, int64:
		value = fmt.Sprintf("%v", val)
	case string:
		value = fmt.Sprintf("%v", val)

		// if opts.ParseEnv is true
		if c.opts.ParseEnv && strings.Index(value, "${") == 0{
			var name, def string
			str := strings.Trim(strings.TrimSpace(value), "${}")
			ss := strings.SplitN(str, "|", 2)

			// ${NotExist|defValue}
			if len(ss) == 2 {
				name, def = strings.TrimSpace(ss[0]), strings.TrimSpace(ss[1])
				// ${SHELL}
			} else {
				name = ss[0]
			}

			value = os.Getenv(name)
			if value == "" {
				value = def
			}
		}
	default:
		ok = false
	}

	return
}

// DefString get a string value, if not found return default value
func (c *Config) DefString(key string, def string) string {
	if value, ok := c.GetString(key); ok {
		return value
	}

	return def
}

// GetInt
func (c *Config) GetInt(key string) (value int, ok bool) {
	rawVal, ok := c.GetString(key)
	if !ok {
		return
	}

	if value, err := strconv.Atoi(rawVal); err != nil {
		return value, true
	}

	return
}

// DefInt get a int value, if not found return default value
func (c *Config) DefInt(key string, def int) int {
	if value, ok := c.GetInt(key); ok {
		return value
	}

	return def
}

// GetBool Looks up a value for a key in this section and attempts to parse that value as a boolean,
// along with a boolean result similar to a map lookup.
// of following( case insensitive):
//  - true
//  - yes
//  - false
//  - no
//  - 1
//  - 0
// The `ok` boolean will be false in the event that the value could not be parsed as a bool
func (c *Config) GetBool(key string) (value bool, ok bool) {
	rawVal, ok := c.GetString(key)
	if !ok {
		return
	}

	lowerCase := strings.ToLower(rawVal)
	switch lowerCase {
	case "", "0", "false", "no":
		value = false
	case "1", "true", "yes":
		value = true
	default:
		ok = false
	}

	return
}

// DefBool get a bool value, if not found return default value
func (c *Config) DefBool(key string, def bool) bool {
	if value, ok := c.GetBool(key); ok {
		return value
	}

	return def
}

// GetStringArr  get config data as a slice/array
func (c *Config) GetStringArr(key string) (arr []string, ok bool) {
	rawVal, ok := c.Get(key)
	if !ok {
		return
	}

	switch rawVal.(type) {
	case []interface{}:
		for _, v := range rawVal.([]interface{}) {
			arr = append(arr, fmt.Sprintf("%v", v))
		}
	default:
		ok = false
	}

	return
}

// GetStringMap get config data as a map[string]string
func (c *Config) GetStringMap(key string) (mp map[string]string, ok bool) {
	rawVal, ok := c.Get(key)
	if !ok {
		return
	}

	switch rawVal.(type) {
	case map[interface{}]interface{}:
		// init map
		mp = make(map[string]string)
		for k, v := range rawVal.(map[interface{}]interface{}) {
			sk := fmt.Sprintf("%v", k)
			mp[sk] = fmt.Sprintf("%v", v)
		}
	default:
		ok = false
	}

	return
}

// MapStructure alias method of the 'GetStructure'
func (c *Config) MapStructure(key string, v interface{}) (err error) {
	return c.GetStructure(key, v)
}

// GetStructure get config data and map to a structure.
// usage:
// 	dbInfo := Db{}
// 	config.GetStructure("db", &dbInfo)
func (c *Config) GetStructure(key string, v interface{}) (err error) {
	if rawVal, ok := c.Get(key); ok {
		blob, err := json.Marshal(rawVal)
		if err != nil {
			return err
		}

		err = json.Unmarshal(blob, v)
	}

	return
}
