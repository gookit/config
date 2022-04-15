package config

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetData(t *testing.T) {
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
	st := assert.New(t)

	c := Default()

	// clear old
	ClearAll()
	// err := LoadFiles("testdata/json_base.json")
	err := LoadStrings(JSON, jsonStr)
	st.Nil(err)

	val := String("name")
	st.Equal("app", val)

	// empty key
	err = Set("", "val")
	st.Error(err)

	// set new value: int
	err = Set("newInt", 23)
	if st.Nil(err) {
		iv := Int("newInt")
		st.Equal(23, iv)
	}

	// set new value: int
	err = Set("newBool", false)
	if st.Nil(err) {
		bv := Bool("newBool")
		st.False(bv)
	}

	// set new value: string
	err = Set("newKey", "new val")
	if st.Nil(err) {
		val = String("newKey")
		st.Equal("new val", val)
	}

	// like yml decoded data
	err = Set("ymlLike", map[interface{}]interface{}{"k": "v"})
	st.Nil(err)
	str := c.String("ymlLike.k")
	st.Equal("v", str)

	err = Set("ymlLike.nk", "nv")
	st.Nil(err)
	str = c.String("ymlLike.nk")
	st.Equal("nv", str)

	// disable setByPath
	err = Set("some.key", "val", false)
	if st.Nil(err) {
		val = String("some")
		st.Equal("", val)

		val = String("some.key")
		st.Equal("val", val)
	}
	// fmt.Printf("%#v\n", c.Data())

	// set value
	err = Set("name", "new name")
	if st.Nil(err) {
		val = String("name")
		st.Equal("new name", val)
	}

	// set value to arr: by path
	err = Set("arr1.1", "new val")
	if st.Nil(err) {
		val = String("arr1.1")
		st.Equal("new val", val)
	}

	// array only support add 1 level value
	err = Set("arr1.1.key", "new val")
	st.Error(err)

	// set value to map: by path
	err = Set("map1.key", "new val")
	if st.Nil(err) {
		val = String("map1.key")
		st.Equal("new val", val)
	}

	// more path nodes
	err = Set("map1.info.key", "val200")
	if st.Nil(err) {
		// fmt.Printf("%v\n", c.Data())
		smp := StringMap("map1.info")
		st.Equal("val200", smp["key"])

		str = String("map1.info.key")
		st.Equal("val200", str)
	}

	// new map
	err = Set("map2.key", "new val")
	if st.Nil(err) {
		val = String("map2.key")

		st.Equal("new val", val)
	}

	// set new value: array(slice)
	err = Set("newArr", []string{"a", "b"})
	if st.Nil(err) {
		arr := Strings("newArr")

		st.Equal(`[]string{"a", "b"}`, fmt.Sprintf("%#v", arr))

		val = String("newArr.1")
		st.Equal("b", val)

		val = String("newArr.100")
		st.Equal("", val)
	}

	// set new value: map
	err = Set("newMap", map[string]string{"k1": "a", "k2": "b"})
	if st.Nil(err) {
		mp := StringMap("newMap")
		st.NotEmpty(mp)
		// st.Equal("map[k1:a k2:b]", fmt.Sprintf("%v", mp))

		val = String("newMap.k1")
		st.Equal("a", val)

		val = String("newMap.notExist")
		st.Equal("", val)
	}

	st.NoError(Set("name.sub", []int{2}))
	ints := Ints("name.sub")
	st.Equal([]int{2}, ints)

	// Readonly
	Default().Readonly()
	st.True(c.Options().Readonly)
	st.Error(Set("name", "new name"))
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
