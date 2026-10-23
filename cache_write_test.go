package config

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

func TestSetInvalidatesCaches(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		c := New("test", EnableCache)
		c.SetData(map[string]any{"name": "old"})
		assert.Eq(t, "old", c.String("name"))
		assert.NoErr(t, c.Set("name", "new"))
		assert.Eq(t, "new", c.String("name"))
	})

	t.Run("replaced parent", func(t *testing.T) {
		c := New("test", EnableCache)
		c.SetData(map[string]any{"profile": map[string]any{"name": "old"}})
		assert.Eq(t, "old", c.String("profile.name"))
		assert.NoErr(t, c.Set("profile", map[string]any{"name": "new"}))
		assert.Eq(t, "new", c.String("profile.name"))
	})

	t.Run("string slice", func(t *testing.T) {
		c := New("test", EnableCache)
		c.SetData(map[string]any{"names": []any{"old"}})
		assert.Eq(t, []string{"old"}, c.Strings("names"))
		assert.NoErr(t, c.Set("names", []any{"new"}))
		assert.Eq(t, []string{"new"}, c.Strings("names"))
	})

	t.Run("updated map child", func(t *testing.T) {
		c := New("test", EnableCache)
		c.SetData(map[string]any{"profile": map[string]any{"name": "old"}})
		assert.Eq(t, map[string]string{"name": "old"}, c.StringMap("profile"))
		assert.NoErr(t, c.Set("profile.name", "new"))
		assert.Eq(t, map[string]string{"name": "new"}, c.StringMap("profile"))
	})

	t.Run("hook", func(t *testing.T) {
		var seen string
		c := New("test", EnableCache, WithHookFunc(func(event string, c *Config) {
			if event == OnSetValue {
				seen = c.String("name")
			}
		}))
		c.SetData(map[string]any{"name": "old"})
		assert.Eq(t, "old", c.String("name"))
		assert.NoErr(t, c.Set("name", "new"))
		assert.Eq(t, "new", seen)
	})
}

func TestSetDataInvalidatesCaches(t *testing.T) {
	var seen string
	c := New("test", EnableCache, WithHookFunc(func(event string, c *Config) {
		if event == OnSetData {
			seen = c.String("name")
		}
	}))
	c.SetData(map[string]any{
		"name":    "old",
		"names":   []any{"old"},
		"profile": map[string]any{"name": "old"},
	})
	assert.Eq(t, "old", c.String("name"))
	assert.Eq(t, []string{"old"}, c.Strings("names"))
	assert.Eq(t, map[string]string{"name": "old"}, c.StringMap("profile"))
	c.SetData(map[string]any{
		"name":    "new",
		"names":   []any{"new"},
		"profile": map[string]any{"name": "new"},
	})
	assert.Eq(t, "new", seen)
	assert.Eq(t, "new", c.String("name"))
	assert.Eq(t, []string{"new"}, c.Strings("names"))
	assert.Eq(t, map[string]string{"name": "new"}, c.StringMap("profile"))

	c.SetData(map[string]any{})
	assert.Eq(t, "", seen)
	assert.Eq(t, "missing", c.String("name", "missing"))
	assert.Empty(t, c.Strings("names"))
	assert.Empty(t, c.StringMap("profile"))
}
