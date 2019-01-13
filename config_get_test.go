package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	st := assert.New(t)

	ClearAll()
	err := LoadStrings(JSON, jsonStr)
	st.Nil(err)

	// fmt.Printf("%#v\n", Data())
	c := Default()

	st.False(c.IsEmpty())

	// error on get
	_, ok := c.Get("")
	st.False(ok)

	_, ok = c.Get("notExist")
	st.False(ok)
	_, ok = c.Get("name.sub")
	st.False(ok)
	st.Error(c.Error())

	_, ok = c.Get("map1.key", false)
	st.False(ok)

	val, ok := Get("map1.notExist")
	st.False(ok)
	st.Nil(val)

	val, ok = Get("notExist.sub")
	st.False(ok)
	st.Nil(val)

	val, ok = c.Get("arr1.100")
	st.False(ok)
	st.Nil(val)

	val, ok = c.Get("arr1.notExist")
	st.False(ok)
	st.Nil(val)

	// get int
	val, ok = Get("age")
	st.True(ok)
	st.Equal(float64(123), val)
	st.Equal("float64", fmt.Sprintf("%T", val))

	iv, ok := Int("age")
	st.True(ok)
	st.Equal(123, iv)

	iv, ok = Int("name")
	st.False(ok)

	iv = DefInt("notExist", 34)
	st.Equal(34, iv)

	iv = c.MustInt("age")
	st.Equal(123, iv)
	iv = c.MustInt("notExist")
	st.Equal(0, iv)

	// get int64
	iv64, ok := Int64("age")
	st.True(ok)
	st.Equal(int64(123), iv64)

	_, ok = Int64("name")
	st.False(false)

	iv64 = DefInt64("age", 34)
	st.Equal(int64(123), iv64)
	iv64 = DefInt64("notExist", 34)
	st.Equal(int64(34), iv64)

	iv64 = c.MustInt64("age")
	st.Equal(int64(123), iv64)
	iv64 = c.MustInt64("notExist")
	st.Equal(int64(0), iv64)

	// get bool
	val, ok = Get("debug")
	st.True(ok)
	st.Equal(true, val)

	bv, ok := Bool("debug")
	st.True(ok)
	st.Equal(true, bv)

	bv, ok = Bool("age")
	st.False(ok)
	st.Equal(false, bv)

	bv = DefBool("debug", false)
	st.Equal(true, bv)

	bv = DefBool("notExist", false)
	st.Equal(false, bv)

	bv = c.MustBool("debug")
	st.True(bv)
	bv = c.MustBool("notExist")
	st.False(bv)

	// get string
	val, ok = Get("name")
	st.True(ok)
	st.Equal("app", val)

	val, ok = String("arr1")
	st.False(ok)
	st.Equal("", val)

	str, ok := String("notExists")
	st.False(ok)
	st.Equal("", str)

	str = DefString("notExists", "defVal")
	st.Equal("defVal", str)

	str = c.MustString("name")
	st.Equal("app", str)
	str = c.MustString("notExist")
	st.Equal("", str)

	// get float
	err = c.Set("flVal", 23.45)
	st.Nil(err)
	flt, ok := c.Float("flVal")
	st.True(ok)
	st.Equal(23.45, flt)

	flt, ok = Float("name")
	st.False(ok)
	st.Equal(float64(0), flt)

	flt = c.DefFloat("notExists", 0)
	st.Equal(float64(0), flt)

	flt = DefFloat("flVal", 0)
	st.Equal(23.45, flt)

	// get string array
	arr, ok := Strings("notExist")
	st.False(ok)

	arr, ok = Strings("map1")
	st.False(ok)

	arr, ok = Strings("arr1")
	st.True(ok)
	st.Equal(`[]string{"val", "val1", "val2"}`, fmt.Sprintf("%#v", arr))

	val, ok = String("arr1.1")
	st.True(ok)
	st.Equal("val1", val)

	err = LoadStrings(JSON, `{
"iArr": [12, 34, 36],
"iMap": {"k1": 12, "k2": 34, "k3": 36}
}`)
	st.Nil(err)

	// Ints: get int arr
	iarr, ok := Ints("name")
	st.False(ok)

	iarr, ok = Ints("notExist")
	st.False(ok)

	iarr, ok = Ints("iArr")
	st.True(ok)
	st.Equal(`[]int{12, 34, 36}`, fmt.Sprintf("%#v", iarr))

	iv, ok = Int("iArr.1")
	st.True(ok)
	st.Equal(34, iv)

	iv, ok = Int("iArr.100")
	st.False(ok)

	// IntMap: get int map
	imp, ok := IntMap("name")
	st.False(ok)
	imp, ok = IntMap("notExist")
	st.False(ok)

	imp, ok = IntMap("iMap")
	st.True(ok)
	st.NotEmpty(imp)

	iv, ok = Int("iMap.k2")
	st.True(ok)
	st.Equal(34, iv)

	iv, ok = Int("iMap.notExist")
	st.False(ok)

	// set a intMap
	err = Set("intMap0", map[string]int{"a": 1, "b": 2})
	st.Nil(err)
	imp, ok = IntMap("intMap0")
	st.True(ok)
	st.NotEmpty(imp)
	st.Equal(1, imp["a"])

	// StringMap: get string map
	smp, ok := StringMap("map1")
	st.True(ok)
	st.Equal("val1", smp["key1"])

	// like load from yaml content
	// c = New("test")
	err = c.LoadData(map[string]interface{}{
		"newIArr":    []int{2, 3},
		"newSArr":    []string{"a", "b"},
		"newIArr1":   []interface{}{12, 23},
		"newIArr2":   []interface{}{12, "abc"},
		"invalidMap": map[string]int{"k": 1},
		"yMap": map[interface{}]interface{}{
			"k0": "v0",
			"k1": 23,
		},
		"yMap1": map[interface{}]interface{}{
			"k":  "v",
			"k1": 23,
			"k2": []interface{}{23, 45},
		},
		"yMap10": map[string]interface{}{
			"k":  "v",
			"k1": 23,
			"k2": []interface{}{23, 45},
		},
		"yMap2": map[interface{}]interface{}{
			"k":  2,
			"k1": 23,
		},
		"yArr": []interface{}{23, 45, "val", map[string]interface{}{"k4": "v4"}},
	})
	st.Nil(err)

	iarr, ok = Ints("newIArr")
	st.True(ok)
	st.Equal("[2 3]", fmt.Sprintf("%v", iarr))

	iarr, ok = Ints("newIArr1")
	st.True(ok)
	st.Equal("[12 23]", fmt.Sprintf("%v", iarr))
	iarr, ok = Ints("newIArr2")
	st.False(ok)

	iv, ok = Int("newIArr.1")
	st.True(ok)
	st.Equal(3, iv)

	iv, ok = Int("newIArr.200")
	st.False(ok)

	// invalid intMap
	imp, ok = IntMap("yMap1")
	st.False(ok)

	imp, ok = IntMap("yMap10")
	st.False(ok)

	imp, ok = IntMap("yMap2")
	st.True(ok)
	st.Equal(2, imp["k"])

	val, ok = String("newSArr.1")
	st.True(ok)
	st.Equal("b", val)

	val, ok = String("newSArr.100")
	st.False(ok)

	smp, ok = StringMap("invalidMap")
	st.False(ok)

	smp, ok = StringMap("yMap.notExist")
	st.False(ok)

	smp, ok = StringMap("yMap")
	st.True(ok)
	st.Equal("v0", smp["k0"])

	iarr, ok = Ints("yMap1.k2")
	st.True(ok)
	st.Equal("[23 45]", fmt.Sprintf("%v", iarr))
}

func TestConfigGetWithDefault(t *testing.T) {
	ClearAll()
	err := LoadStrings(JSON, jsonStr)
	assert.Nil(t, err)

	// fmt.Printf("%#v\n", Data())
	c := Default()
	assert.Equal(t, 0, c.DefInt("not-exist"))
	assert.Equal(t, int64(0), c.DefInt64("not-exist"))
	assert.Equal(t, float64(0), c.DefFloat("not-exist"))
	assert.Equal(t, false, c.DefBool("not-exist"))
	assert.Equal(t, "", c.DefString("not-exist"))
}

func TestConfig_MapStructure(t *testing.T) {
	st := assert.New(t)

	cfg := New("test")
	err := cfg.LoadStrings(JSON, `{
"age": 28,
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)

	st.Nil(err)

	user := &struct {
		Age    int
		Name   string
		Sports []string
	}{}
	// map all
	err = cfg.MapStruct("", user)
	st.Nil(err)

	st.Equal(28, user.Age)
	st.Equal("inhere", user.Name)
	st.Equal("pingPong", user.Sports[0])

	// map some
	err = cfg.LoadStrings(JSON, `{
"sec": {
	"key": "val",
	"age": 120,
	"tags": [12, 34]
}
}`)
	st.Nil(err)

	some := struct {
		Age  int
		Kye  string
		Tags []int
	}{}
	err = cfg.MapStructure("sec", &some)
	st.Nil(err)
	st.Equal(120, some.Age)
	st.Equal(12, some.Tags[0])
}
