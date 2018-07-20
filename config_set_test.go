package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {
	st := assert.New(t)

	// clear old
	ClearAll()
	// err := LoadFiles("testdata/json_base.json")
	err := LoadStrings(Json, jsonStr)
	st.Nil(err)

	val, ok := String("name")
	st.True(ok)
	st.Equal("app", val)

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

	// set value to map: by path
	err = Set("map1.key", "new val")
	if st.Nil(err) {
		val, ok = String("map1.key")
		st.True(ok)
		st.Equal("new val", val)
	}

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

	// set new value: array(slice)
	err = Set("newArr", []string{"a", "b"})
	if st.Nil(err) {
		arr, ok := Strings("newArr")
		st.True(ok)
		st.Equal(`[]string{"a", "b"}`, fmt.Sprintf("%#v", arr))

		val, ok = String("newArr.1")
		st.True(ok)
		st.Equal("b", val)
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
	}
}
