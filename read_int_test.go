package config

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/x/assert"
)

// Note: on 32-bit platforms(eg GOARCH=386) the strconv bitSize 0 is the platform
// int width, reading the full int64 range must not be clamped to the int32 range.
// See https://github.com/gookit/config/pull/228

func TestInt64FullRange(t *testing.T) {
	for _, value := range []int64{1 << 31, -1<<31 - 1, 1<<63 - 1, -1 << 63} {
		t.Run(fmt.Sprint(value), func(t *testing.T) {
			is := assert.New(t)
			c := New("int64-range")

			is.NoErr(c.Set("value", value))
			is.Eq(value, c.Int64("value"))
			is.NoErr(c.Error())
		})
	}
}

// TestIntsIntMapParseEnv: Ints and IntMap convert a "${VAR}" elem/value after
// ParseEnvValue when ParseEnv is enabled, like Strings and StringMap do.
func TestIntsIntMapParseEnv(t *testing.T) {
	testutil.MockEnvValue("CONFIG_READ_INT_PORT", "8080", func(_ string) {
		is := assert.New(t)

		c := New("ints-env", ParseEnv)
		is.NoErr(c.Set("list", []any{"${CONFIG_READ_INT_PORT}", 1, "2"}))
		is.NoErr(c.Set("strs", []string{"${CONFIG_READ_INT_PORT}", "9"}))
		is.NoErr(c.Set("smap", map[string]string{"http": "${CONFIG_READ_INT_PORT}"}))
		is.NoErr(c.Set("amap", map[string]any{"http": "${CONFIG_READ_INT_PORT}"}))

		is.Eq([]int{8080, 1, 2}, c.Ints("list"))
		is.Eq([]int{8080, 9}, c.Ints("strs"))
		is.Eq(map[string]int{"http": 8080}, c.IntMap("smap"))
		is.Eq(map[string]int{"http": 8080}, c.IntMap("amap"))
		is.NoErr(c.Error())

		// ParseEnv off: no env parsing, the numbers are converted as before
		c = New("ints-no-env")
		is.NoErr(c.Set("list", []any{1, "2"}))
		is.NoErr(c.Set("strs", []string{"3"}))
		is.Eq([]int{1, 2}, c.Ints("list"))
		is.Eq([]int{3}, c.Ints("strs"))
		is.NoErr(c.Error())
	})
}

func TestInt64OutOfRange(t *testing.T) {
	// like Uint, an out of range value records an error and the clamped value is not returned
	for _, value := range []string{"9223372036854775808", "-9223372036854775809"} {
		t.Run(value, func(t *testing.T) {
			is := assert.New(t)
			c := New("int64-out-of-range")

			is.NoErr(c.Set("value", value))
			is.Eq(int64(0), c.Int64("value"))
			is.Err(c.Error())
			is.Eq(0, c.Int("value"))
		})
	}
}

func TestInt64ByDuration(t *testing.T) {
	is := assert.New(t)
	c := New("duration-range")

	// time.Duration is an int64 type, it must keep the full int64 range on 32-bit platforms
	is.NoErr(c.Set("value", int64(1)<<40))
	is.Eq(time.Duration(1)<<40, c.Duration("value"))
	is.NoErr(c.Error())
}

func TestIntWithInt64Value(t *testing.T) {
	is := assert.New(t)
	c := New("int-range")

	value := int64(1)<<40 + 3
	is.NoErr(c.Set("value", value))

	got := c.Int("value")
	if strconv.IntSize == 64 {
		is.Eq(int(value), got)
		is.NoErr(c.Error())
	} else {
		// out of the int range on 32-bit platforms
		is.Eq(0, got)
		is.Err(c.Error())
	}
}

func TestUintWithInt64Value(t *testing.T) {
	is := assert.New(t)

	// negative numbers are not a valid uint value
	c := New("uint-range")
	is.NoErr(c.Set("value", int64(-1)))
	is.Eq(uint(0), c.Uint("value"))
	is.Err(c.Error())

	value := int64(1)<<40 + 3
	c2 := New("uint-range2")
	is.NoErr(c2.Set("value", value))

	got := c2.Uint("value")
	if strconv.IntSize == 64 {
		is.Eq(uint(value), got)
		is.NoErr(c2.Error())
	} else {
		// out of the uint range on 32-bit platforms
		is.Eq(uint(0), got)
		is.Err(c2.Error())
	}
}
