package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	st := assert.New(t)

	ClearAll()
	err := LoadStrings(Json, jsonStr)
	st.Nil(err)

	// fmt.Printf("%#v\n", Data())

	// get int
	iv, ok := Get("age")
	st.True(ok)
	st.Equal(float64(123), iv)
	st.Equal("float64", fmt.Sprintf("%T", iv))

	iv, ok = Int("age")
	st.True(ok)
	st.Equal(123, iv)

	iv = DefInt("notExist", 34)
	st.Equal(34, iv)

	// get int64
	iv64, ok := Int64("age")
	st.True(ok)
	st.Equal(int64(123), iv64)

	iv64 = DefInt64("notExist", 34)
	st.Equal(int64(34), iv64)

	// get bool
	bv, ok := Get("debug")
	st.True(ok)
	st.Equal(true, bv)

	bv, ok = Bool("debug")
	st.True(ok)
	st.Equal(true, bv)

	bv = DefBool("notExist", false)
	st.Equal(false, bv)

	// get string
	val, ok := Get("name")
	st.True(ok)
	st.Equal("app", val)

	// get not exists
	val, ok = Get("notExist", false)
	st.False(ok)
	st.Nil(val)

	str, ok := String("notExists")
	st.False(ok)
	st.Equal("", str)

	def := DefString("notExists", "defVal")
	st.Equal("defVal", def)

	// get array
	arr, ok := Strings("arr1")
	st.True(ok)
	st.Equal(`[]string{"val", "val1", "val2"}`, fmt.Sprintf("%#v", arr))

	val, ok = String("arr1.1")
	st.True(ok)
	st.Equal("val1", val)

	err = LoadStrings(Json, `{
"iArr": [12, 34, 36],
"iMap": {"k1": 12, "k2": 34, "k3": 36}
}`)
	st.Nil(err)

	// get int arr
	iarr, ok := Ints("iArr")
	st.True(ok)
	st.Equal(`[]int{12, 34, 36}`, fmt.Sprintf("%#v", iarr))

	iv, ok = Int("iArr.1")
	st.True(ok)
	st.Equal(34, iv)

	// get int map
	imp, ok := IntMap("iMap")
	st.True(ok)
	st.NotEmpty(imp)

	iv, ok = Int("iMap.k2")
	st.True(ok)
	st.Equal(34, iv)
}

type user struct {
	Age    int
	Name   string
	Sports []string
}

func TestConfig_MapStructure(t *testing.T) {
	st := assert.New(t)

	cfg := New("test")
	err := cfg.LoadStrings(Json, `{
"age": 28,
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)

	st.Nil(err)

	user := &user{}
	err = cfg.MapStructure("", user)
	st.Nil(err)

	st.Equal(28, user.Age)
	st.Equal("inhere", user.Name)
	st.Equal("pingPong", user.Sports[0])
}
