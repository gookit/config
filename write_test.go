package config

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func TestSetData(t *testing.T) {
	defer func() {
		Reset()
	}()

	c := Default()

	err := c.LoadStrings(JSON, jsonStr)
	assert.NoErr(t, err)
	assert.Eq(t, "app", c.String("name"))
	assert.True(t, c.Exists("age"))

	SetData(map[string]any{
		"name": "new app",
	})
	assert.Eq(t, "new app", c.String("name"))
	assert.False(t, c.Exists("age"))

	c.SetData(map[string]any{
		"age": 222,
	})
	assert.Eq(t, "", c.String("name"))
	assert.False(t, c.Exists("name"))
	assert.True(t, c.Exists("age"))
	assert.Eq(t, 222, c.Int("age"))
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
	is.Eq("app", val)

	// empty key
	err = Set("", "val")
	is.Err(err)

	// set new value: int
	err = Set("newInt", 23)
	if is.Nil(err).IsOk() {
		iv := Int("newInt")
		is.Eq(23, iv)
	}

	// set new value: int
	err = Set("newBool", false)
	if is.Nil(err).IsOk() {
		bv := Bool("newBool")
		is.False(bv)
	}

	// set new value: string
	err = Set("newKey", "new val")
	if is.Nil(err).IsOk() {
		val = String("newKey")
		is.Eq("new val", val)
	}

	// like yaml.v2 decoded data
	err = Set("ymlLike", map[any]any{"k": "v"})
	is.Nil(err)
	str := c.String("ymlLike.k")
	is.Eq("v", str)

	err = Set("ymlLike.nk", "nv")
	is.Nil(err)
	str = c.String("ymlLike.nk")
	is.Eq("nv", str)

	// disable setByPath
	err = Set("some.key", "val", false)
	if is.Nil(err).IsOk() {
		val = String("some")
		is.Eq("", val)

		val = String("some.key")
		is.Eq("val", val)
	}
	// fmt.Printf("%#v\n", c.Data())

	// set value
	err = Set("name", "new name")
	if is.Nil(err).IsOk() {
		val = String("name")
		is.Eq("new name", val)
	}

	// set value to arr: by path
	err = Set("arr1.1", "new val")
	if is.Nil(err).IsOk() {
		val = String("arr1.1")
		is.Eq("new val", val)
	}

	// array only support add 1 level value
	err = Set("arr1.1.key", "new val")
	is.Err(err)

	// set value to map: by path
	err = Set("map1.key", "new val")
	if is.Nil(err).IsOk() {
		val = String("map1.key")
		is.Eq("new val", val)
	}

	// more path nodes
	err = Set("map1.info.key", "val200")
	if is.Nil(err).IsOk() {
		// fmt.Printf("%v\n", c.Data())
		smp := StringMap("map1.info")
		is.Eq("val200", smp["key"])

		str = String("map1.info.key")
		is.Eq("val200", str)
	}

	// new map
	err = Set("map2.key", "new val")
	if is.Nil(err).IsOk() {
		val = String("map2.key")

		is.Eq("new val", val)
	}

	// set new value: array(slice)
	err = Set("newArr", []string{"a", "b"})
	if is.Nil(err).IsOk() {
		arr := Strings("newArr")

		is.Eq(`[]string{"a", "b"}`, fmt.Sprintf("%#v", arr))

		val = String("newArr.1")
		is.Eq("b", val)

		val = String("newArr.100")
		is.Eq("", val)
	}

	// set new value: map
	err = Set("newMap", map[string]string{"k1": "a", "k2": "b"})
	if is.Nil(err).IsOk() {
		mp := StringMap("newMap")
		is.NotEmpty(mp)
		// is.Eq("map[k1:a k2:b]", fmt.Sprintf("%v", mp))

		val = String("newMap.k1")
		is.Eq("a", val)

		val = String("newMap.notExist")
		is.Eq("", val)
	}

	is.NoErr(Set("name.sub", []int{2}, false))
	ints := Ints("name.sub")
	is.Eq([]int{2}, ints)

	// Readonly
	Default().Readonly()
	is.True(c.Options().Readonly)
	is.Err(Set("name", "new name"))
}

func TestSet_fireEvent(t *testing.T) {
	var buf bytes.Buffer
	hookFn := func(event string, c *Config) {
		buf.WriteString("fire the: ")
		buf.WriteString(event)
	}

	c := NewWithOptions("test", WithHookFunc(hookFn))
	err := c.LoadData(map[string]any{
		"key": "value",
	})
	assert.NoErr(t, err)
	assert.Eq(t, "fire the: load.data", buf.String())
	buf.Reset()

	err = c.Set("key", "value2")
	assert.NoErr(t, err)
	assert.Eq(t, "fire the: set.value", buf.String())
	buf.Reset()
}
