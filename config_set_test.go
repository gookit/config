package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {
	st := assert.New(t)

	c := Default()

	// clear old
	ClearAll()
	// err := LoadFiles("testdata/json_base.json")
	err := LoadStrings(JSON, jsonStr)
	st.Nil(err)

	val, ok := String("name")
	st.True(ok)
	st.Equal("app", val)

	// empty key
	err = Set("", "val")
	st.Error(err)

	// set new value: int
	err = Set("newInt", 23)
	if st.Nil(err) {
		iv, ok := Int("newInt")
		st.True(ok)
		st.Equal(23, iv)
	}

	// set new value: int
	err = Set("newBool", false)
	if st.Nil(err) {
		bv, ok := Bool("newBool")
		st.True(ok)
		st.False(bv)
	}

	// set new value: string
	err = Set("newKey", "new val")
	if st.Nil(err) {
		val, ok = String("newKey")
		st.True(ok)
		st.Equal("new val", val)
	}

	// like yml decoded data
	err = Set("ymlLike", map[interface{}]interface{}{"k": "v"})
	st.Nil(err)
	str := c.MustString("ymlLike.k")
	st.Equal("v", str)

	err = Set("ymlLike.nk", "nv")
	st.Nil(err)
	str = c.MustString("ymlLike.nk")
	st.Equal("nv", str)

	// disable setByPath
	err = Set("some.key", "val", false)
	if st.Nil(err) {
		val, ok = String("some")
		st.False(ok)
		st.Equal("", val)

		val, ok = String("some.key")
		st.True(ok)
		st.Equal("val", val)
	}
	// fmt.Printf("%#v\n", c.Data())

	// set value
	err = Set("name", "new name")
	if st.Nil(err) {
		val, ok = String("name")
		st.True(ok)
		st.Equal("new name", val)
	}

	// set value to arr: by path
	err = Set("arr1.1", "new val")
	if st.Nil(err) {
		val, ok = String("arr1.1")
		st.True(ok)
		st.Equal("new val", val)
	}

	// array only support add 1 level value
	err = Set("arr1.1.key", "new val")
	st.Error(err)

	// set value to map: by path
	err = Set("map1.key", "new val")
	if st.Nil(err) {
		val, ok = String("map1.key")
		st.True(ok)
		st.Equal("new val", val)
	}

	// more path nodes
	err = Set("map1.info.key", "val200")
	if st.Nil(err) {
		// fmt.Printf("%v\n", c.Data())
		smp, ok := StringMap("map1.info")
		st.True(ok)
		st.Equal("val200", smp["key"])

		str, ok = String("map1.info.key")
		st.True(ok)
		st.Equal("val200", str)
	}

	// new map
	err = Set("map2.key", "new val")
	if st.Nil(err) {
		val, ok = String("map2.key")
		st.True(ok)
		st.Equal("new val", val)
	}

	// set new value: array(slice)
	err = Set("newArr", []string{"a", "b"})
	if st.Nil(err) {
		arr, ok := Strings("newArr")
		st.True(ok)
		st.Equal(`[]string{"a", "b"}`, fmt.Sprintf("%#v", arr))

		val, ok = String("newArr.1")
		st.True(ok)
		st.Equal("b", val)

		val, ok = String("newArr.100")
		st.False(ok)
		st.Equal("", val)
	}

	// set new value: map
	err = Set("newMap", map[string]string{"k1": "a", "k2": "b"})
	if st.Nil(err) {
		mp, ok := StringMap("newMap")
		st.True(ok)
		st.NotEmpty(mp)
		// st.Equal("map[k1:a k2:b]", fmt.Sprintf("%v", mp))

		val, ok = String("newMap.k1")
		st.True(ok)
		st.Equal("a", val)

		val, ok = String("newMap.notExist")
		st.False(ok)
		st.Equal("", val)
	}

	st.NoError(Set("name.sub", []int{2}))
	ints, ok := Ints("name.sub")
	st.True(ok)
	st.Equal([]int{2}, ints)

	// Readonly
	Default().Readonly()
	st.Error(Set("name", "new name"))
}
