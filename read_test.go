package config

import (
	"fmt"
	"testing"

	"github.com/gookit/goutil/testutil"
	"github.com/stretchr/testify/assert"
)

func TestConfig_GetValue(t *testing.T) {
	is := assert.New(t)

	ClearAll()
	err := LoadStrings(JSON, jsonStr)
	is.Nil(err)

	c := Default()

	// error on get
	_, ok := GetValue("")
	is.False(ok)

	_, ok = c.GetValue("notExist")
	is.False(ok)
	_, ok = c.GetValue("name.sub")
	is.False(ok)
	is.Error(c.Error())

	_, ok = c.GetValue("map1.key", false)
	is.False(ok)
	is.False(Exists("map1.key", false))

	val, ok := GetValue("map1.notExist")
	is.Nil(val)
	is.False(ok)
	is.False(Exists("map1.notExist"))

	val, ok = GetValue("notExist.sub")
	is.False(ok)
	is.Nil(val)
	is.False(Exists("notExist.sub"))

	val, ok = c.GetValue("arr1.100")
	is.Nil(val)
	is.False(ok)
	is.False(Exists("arr1.100"))

	val, ok = c.GetValue("arr1.notExist")
	is.Nil(val)
	is.False(ok)
	is.False(Exists("arr1.notExist"))

	// load data for tests
	err = c.LoadData(map[string]interface{}{
		"setStrMap": map[string]string{
			"k": "v",
		},
		"setIntMap": map[string]int{
			"k2": 23,
		},
	})
	is.Nil(err)
	// -- assert map[string]string
	is.True(Exists("setStrMap.k"))
	is.Equal("v", Get("setStrMap.k"))
	is.False(Exists("setStrMap.k1"))

	// -- assert map[string]int
	is.True(Exists("setIntMap.k2"))
	is.Equal(23, Get("setIntMap.k2"))
	is.False(Exists("setIntMap.k1"))

	ClearAll()
}

func TestGet(t *testing.T) {
	is := assert.New(t)

	ClearAll()
	err := LoadStrings(JSON, jsonStr)
	is.Nil(err)

	// fmt.Printf("%#v\n", Data())
	c := Default()

	is.False(c.IsEmpty())
	is.True(Exists("age"))
	is.True(Exists("map1.key"))
	is.True(Exists("arr1.1"))
	is.False(Exists("arr1.1", false))
	is.False(Exists("not-exist.sub"))
	is.False(Exists(""))
	is.False(Exists("not-exist"))

	// get value
	val := Get("age")
	is.Equal(float64(123), val)
	is.Equal("float64", fmt.Sprintf("%T", val))

	val = Get("not-exist")
	is.Nil(val)

	val = Get("name")
	is.Equal("app", val)

	// get string array
	arr := Strings("notExist")
	is.Empty(arr)

	arr = Strings("map1")
	is.Empty(arr)

	arr = Strings("arr1")
	is.Equal(`[]string{"val", "val1", "val2"}`, fmt.Sprintf("%#v", arr))

	val = String("arr1.1")
	is.Equal("val1", val)

	err = LoadStrings(JSON, `{
"iArr": [12, 34, 36],
"iMap": {"k1": 12, "k2": 34, "k3": 36}
}`)
	is.Nil(err)

	// Ints: get int arr
	iarr := Ints("name")
	is.False(Exists("name.1"))
	is.Empty(iarr)

	iarr = Ints("notExist")
	is.Empty(iarr)

	iarr = Ints("iArr")
	is.Equal(`[]int{12, 34, 36}`, fmt.Sprintf("%#v", iarr))

	iv := Int("iArr.1")
	is.Equal(34, iv)

	iv = Int("iArr.100")
	is.Equal(0, iv)

	// IntMap: get int map
	imp := IntMap("name")
	is.Empty(imp)
	imp = IntMap("notExist")
	is.Empty(imp)

	imp = IntMap("iMap")
	is.NotEmpty(imp)

	iv = Int("iMap.k2")
	is.Equal(34, iv)
	is.True(Exists("iMap.k2"))

	iv = Int("iMap.notExist")
	is.Equal(0, iv)
	is.False(Exists("iMap.notExist"))

	// set a intMap
	err = Set("intMap0", map[string]int{"a": 1, "b": 2})
	is.Nil(err)

	imp = IntMap("intMap0")
	is.NotEmpty(imp)
	is.Equal(1, imp["a"])
	is.Equal(2, Get("intMap0.b"))
	is.True(Exists("intMap0.a"))
	is.False(Exists("intMap0.c"))

	// StringMap: get string map
	smp := StringMap("map1")
	is.Equal("val1", smp["key1"])

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
	is.Nil(err)

	iarr = Ints("newIArr")
	is.Equal("[2 3]", fmt.Sprintf("%v", iarr))

	iarr = Ints("newIArr1")
	is.Equal("[12 23]", fmt.Sprintf("%v", iarr))
	iarr = Ints("newIArr2")
	is.Empty(iarr)

	iv = Int("newIArr.1")
	is.True(Exists("newIArr.1"))
	is.Equal(3, iv)

	iv = Int("newIArr.200")
	is.False(Exists("newIArr.200"))
	is.Equal(0, iv)

	// invalid intMap
	imp = IntMap("yMap1")
	is.Empty(imp)

	imp = IntMap("yMap10")
	is.Empty(imp)

	imp = IntMap("yMap2")
	is.Equal(2, imp["k"])

	val = String("newSArr.1")
	is.True(Exists("newSArr.1"))
	is.Equal("b", val)

	val = String("newSArr.100")
	is.False(Exists("newSArr.100"))
	is.Equal("", val)

	smp = StringMap("invalidMap")
	is.Nil(smp)

	smp = StringMap("yMap.notExist")
	is.Nil(smp)

	smp = StringMap("yMap")
	is.True(Exists("yMap.k0"))
	is.False(Exists("yMap.k100"))
	is.Equal("v0", smp["k0"])

	iarr = Ints("yMap1.k2")
	is.Equal("[23 45]", fmt.Sprintf("%v", iarr))
}

func TestInt(t *testing.T) {
	is := assert.New(t)
	ClearAll()
	_ = LoadStrings(JSON, jsonStr)

	is.True(Exists("age"))

	iv := Int("age")
	is.Equal(123, iv)

	iv = Int("name")
	is.Equal(0, iv)

	iv = Int("notExist", 34)
	is.Equal(34, iv)

	c := Default()
	iv = c.Int("age")
	is.Equal(123, iv)
	iv = c.Int("notExist")
	is.Equal(0, iv)

	uiv := Uint("age")
	is.Equal(uint(123), uiv)

	ClearAll()
}

func TestInt64(t *testing.T) {
	is := assert.New(t)
	ClearAll()
	_ = LoadStrings(JSON, jsonStr)

	// get int64
	iv64 := Int64("age")
	is.Equal(int64(123), iv64)

	iv64 = Int64("name")
	is.Equal(iv64, int64(0))

	iv64 = Int64("age", 34)
	is.Equal(int64(123), iv64)
	iv64 = Int64("notExist", 34)
	is.Equal(int64(34), iv64)

	c := Default()
	iv64 = c.Int64("age")
	is.Equal(int64(123), iv64)
	iv64 = c.Int64("notExist")
	is.Equal(int64(0), iv64)

	ClearAll()
}

func TestFloat(t *testing.T) {
	is := assert.New(t)
	ClearAll()
	_ = LoadStrings(JSON, jsonStr)
	c := Default()

	// get float
	err := c.Set("flVal", 23.45)
	is.Nil(err)
	flt := c.Float("flVal")
	is.Equal(23.45, flt)

	flt = Float("name")
	is.Equal(float64(0), flt)

	flt = c.Float("notExists")
	is.Equal(float64(0), flt)

	flt = c.Float("notExists", 10)
	is.Equal(float64(10), flt)

	flt = Float("flVal", 0)
	is.Equal(23.45, flt)

	ClearAll()
}

func TestString(t *testing.T) {
	is := assert.New(t)
	ClearAll()
	_ = LoadStrings(JSON, jsonStr)

	// get string
	val := String("arr1")
	is.Equal("[val val1 val2]", val)

	str := String("notExists")
	is.Equal("", str)

	str = String("notExists", "defVal")
	is.Equal("defVal", str)

	c := Default()
	str = c.String("name")
	is.Equal("app", str)
	str = c.String("notExist")
	is.Equal("", str)

	ClearAll()
}

func TestBool(t *testing.T) {
	is := assert.New(t)
	ClearAll()
	_ = LoadSources(JSON, []byte(jsonStr))

	// get bool
	val := Get("debug")
	is.Equal(true, val)

	bv := Bool("debug")
	is.Equal(true, bv)

	bv = Bool("age")
	is.Equal(false, bv)

	bv = Bool("debug", false)
	is.Equal(true, bv)

	bv = Bool("notExist", false)
	is.Equal(false, bv)

	c := Default()
	bv = c.Bool("debug")
	is.True(bv)
	bv = c.Bool("notExist")
	is.False(bv)

	ClearAll()
}

func TestParseEnv(t *testing.T) {
	is := assert.New(t)

	cfg := NewWithOptions("test", ParseEnv)
	err := cfg.LoadStrings(JSON, `{
"ekey": "${EnvKey}",
"ekey0": "${ EnvKey0 }",
"ekey1": "${EnvKey1|defValue}",
"ekey2": "${ EnvKey2 | defValue1 }",
"ekey3": "${ EnvKey3 | app:run }",
"ekey4": "${FirstEnv}/${ SecondEnv }",
"ekey5": "${TEST_SHELL|/bin/bash}",
"ekey6": "${ EnvKey6 | app=run }",
"ekey7": "${ EnvKey7 | app.run }",
"ekey8": "${ EnvKey8 | app/run }"
}`)

	is.NoError(err)

	tests := []struct{ EKey, EVal, CKey, CVal string }{
		{"EnvKey", "EnvKey val", "ekey", "EnvKey val"},
		{"EnvKey", "", "ekey", "${EnvKey}"},
		{"EnvKey0", "EnvKey0 val", "ekey0", "EnvKey0 val"},
		{"EnvKey3", "EnvKey3 val", "ekey3", "EnvKey3 val"},
		{"EnvKey3", "", "ekey3", "app:run"},
		{"EnvKey6", "", "ekey6", "app=run"},
		{"EnvKey7", "", "ekey7", "app.run"},
		{"EnvKey8", "", "ekey8", "app/run"},
		{"TEST_SHELL", "/bin/zsh", "ekey5", "/bin/zsh"},
		{"TEST_SHELL", "", "ekey5", "/bin/bash"},
	}

	for _, smp := range tests {
		is.Equal("", Getenv(smp.EKey))

		testutil.MockEnvValue(smp.EKey, smp.EVal, func(eVal string) {
			is.Equal(smp.EVal, eVal)
			is.Equal(smp.CVal, cfg.String(smp.CKey))
		})
	}

	// test multi ENV key
	is.Equal("", Getenv("FirstEnv"))

	testutil.MockEnvValues(map[string]string{
		"FirstEnv":  "abc",
		"SecondEnv": "def",
	}, func() {
		is.Equal("abc", Getenv("FirstEnv"))
		is.Equal("def", Getenv("SecondEnv"))
		is.Equal("abc/def", cfg.String("ekey4"))
	})

	testutil.MockEnvValues(map[string]string{
		"FirstEnv": "abc",
	}, func() {
		is.Equal("abc", Getenv("FirstEnv"))
		is.Equal("", Getenv("SecondEnv"))
		is.Equal("abc/${ SecondEnv }", cfg.String("ekey4"))
	})
}
