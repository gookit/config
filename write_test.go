package config

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetData(t *testing.T) {
	defer func() {
		Reset()
	}()

	c := Default()

	err := c.LoadStrings(JSON, jsonStr)
	assert.NoError(t, err)
	assert.Equal(t, "app", c.String("name"))
	assert.True(t, c.Exists("age"))

	SetData(map[string]interface{}{
		"name": "new app",
	})
	assert.Equal(t, "new app", c.String("name"))
	assert.False(t, c.Exists("age"))

	c.SetData(map[string]interface{}{
		"age": 222,
	})
	assert.Equal(t, "", c.String("name"))
	assert.False(t, c.Exists("name"))
	assert.True(t, c.Exists("age"))
	assert.Equal(t, 222, c.Int("age"))
}

func TestSet(t *testing.T) {
	defer func() {
		ClearAll()
	}()

	is := assert.New(t)
	c := Default()

	// clear old
	ClearAll()
	// err := LoadFiles("testdata/json_base.json")
	err := LoadStrings(JSON, jsonStr)
	is.Nil(err)

	val := String("name")
	is.Equal("app", val)

	// empty key
	err = Set("", "val")
	is.Error(err)

	// set new value: int
	err = Set("newInt", 23)
	if is.Nil(err) {
		iv := Int("newInt")
		is.Equal(23, iv)
	}

	// set new value: int
	err = Set("newBool", false)
	if is.Nil(err) {
		bv := Bool("newBool")
		is.False(bv)
	}

	// set new value: string
	err = Set("newKey", "new val")
	if is.Nil(err) {
		val = String("newKey")
		is.Equal("new val", val)
	}

	// like yaml.v2 decoded data
	err = Set("ymlLike", map[interface{}]interface{}{"k": "v"})
	is.Nil(err)
	str := c.String("ymlLike.k")
	is.Equal("v", str)

	err = Set("ymlLike.nk", "nv")
	is.Nil(err)
	str = c.String("ymlLike.nk")
	is.Equal("nv", str)

	// disable setByPath
	err = Set("some.key", "val", false)
	if is.Nil(err) {
		val = String("some")
		is.Equal("", val)

		val = String("some.key")
		is.Equal("val", val)
	}
	// fmt.Printf("%#v\n", c.Data())

	// set value
	err = Set("name", "new name")
	if is.Nil(err) {
		val = String("name")
		is.Equal("new name", val)
	}

	// set value to arr: by path
	err = Set("arr1.1", "new val")
	if is.Nil(err) {
		val = String("arr1.1")
		is.Equal("new val", val)
	}

	// array only support add 1 level value
	err = Set("arr1.1.key", "new val")
	is.Error(err)

	// set value to map: by path
	err = Set("map1.key", "new val")
	if is.Nil(err) {
		val = String("map1.key")
		is.Equal("new val", val)
	}

	// more path nodes
	err = Set("map1.info.key", "val200")
	if is.Nil(err) {
		// fmt.Printf("%v\n", c.Data())
		smp := StringMap("map1.info")
		is.Equal("val200", smp["key"])

		str = String("map1.info.key")
		is.Equal("val200", str)
	}

	// new map
	err = Set("map2.key", "new val")
	if is.Nil(err) {
		val = String("map2.key")

		is.Equal("new val", val)
	}

	// set new value: array(slice)
	err = Set("newArr", []string{"a", "b"})
	if is.Nil(err) {
		arr := Strings("newArr")

		is.Equal(`[]string{"a", "b"}`, fmt.Sprintf("%#v", arr))

		val = String("newArr.1")
		is.Equal("b", val)

		val = String("newArr.100")
		is.Equal("", val)
	}

	// set new value: map
	err = Set("newMap", map[string]string{"k1": "a", "k2": "b"})
	if is.Nil(err) {
		mp := StringMap("newMap")
		is.NotEmpty(mp)
		// is.Equal("map[k1:a k2:b]", fmt.Sprintf("%v", mp))

		val = String("newMap.k1")
		is.Equal("a", val)

		val = String("newMap.notExist")
		is.Equal("", val)
	}

	is.NoError(Set("name.sub", []int{2}, false))
	ints := Ints("name.sub")
	is.Equal([]int{2}, ints)

	// Readonly
	Default().Readonly()
	is.True(c.Options().Readonly)
	is.Error(Set("name", "new name"))
}

func TestSet_fireEvent(t *testing.T) {
	var buf bytes.Buffer
	hookFn := func(event string, c *Config) {
		buf.WriteString("fire the: ")
		buf.WriteString(event)
	}

	c := NewWithOptions("test", WithHookFunc(hookFn))
	err := c.LoadData(map[string]interface{}{
		"key": "value",
	})
	assert.NoError(t, err)
	assert.Equal(t, "fire the: load.data", buf.String())
	buf.Reset()

	err = c.Set("key", "value2")
	assert.NoError(t, err)
	assert.Equal(t, "fire the: set.value", buf.String())
	buf.Reset()
}
