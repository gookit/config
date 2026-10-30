package config_test

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/goutil/x/assert"
	"testing"
)

func TestSetAfterNilData(t *testing.T) {
	cases := []struct {
		name, key string
		parsePath bool
	}{
		{"top level", "host", true},
		{"nested path", "db.host", true},
		{"literal dotted key", "db.host", false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := config.New("nil-data")
			c.SetData(nil)
			defer func() {
				if p := recover(); p != nil {
					t.Fatalf("Set after nil data panicked: %v", p)
				}
			}()
			assert.NoErr(t, c.Set(tt.key, "localhost", tt.parsePath))
			assert.Eq(t, "localhost", c.Get(tt.key, tt.parsePath))
			assert.False(t, c.IsEmpty())
		})
	}
	t.Run("readonly", func(t *testing.T) {
		c := config.New("readonly", config.Readonly)
		c.SetData(nil)
		assert.Eq(t, config.ErrReadonly, c.Set("host", "localhost"))
		assert.Nil(t, c.Data())
	})
	t.Run("invalid key", func(t *testing.T) {
		c := config.New("invalid")
		c.SetData(nil)
		assert.Eq(t, config.ErrKeyIsEmpty, c.Set("", "localhost"))
		assert.Nil(t, c.Data())
	})
}
