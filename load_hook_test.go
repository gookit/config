package config

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gookit/goutil/x/assert"
)

func TestLoadHookSeesNewValue(t *testing.T) {
	var seen string
	c := New("load-hook", EnableCache, WithHookFunc(func(event string, c *Config) {
		if event == OnLoadData {
			seen = c.String("name")
		}
	}))

	assert.NoErr(t, c.LoadData(map[string]any{"name": "old"}))
	assert.Eq(t, "old", seen)

	// the cache is filled by the hook read above, it must not hide the new value
	assert.NoErr(t, c.LoadData(map[string]any{"name": "new"}))
	assert.Eq(t, "new", seen)
}

// the hook is fired without the data lock, so another writer must not change the data
// while the hook is running.
func TestLoadHookKeepsDataStable(t *testing.T) {
	is := assert.New(t)

	hookEnter := make(chan struct{})
	hookRelease := make(chan struct{})
	var firstHook atomic.Bool

	c := New("load-stable", WithHookFunc(func(event string, cfg *Config) {
		// hold only the first hook until the test allows it to return
		if event == OnLoadData && firstHook.CompareAndSwap(false, true) {
			close(hookEnter)
			<-hookRelease
		}
	}))

	go func() { _ = c.LoadData(map[string]any{"first": 1}) }()
	<-hookEnter

	// another writer must wait for the running hook
	done := make(chan struct{})
	go func() {
		_ = c.Set("second", 2)
		close(done)
	}()

	select {
	case <-done:
		close(hookRelease)
		t.Fatal("another Set completed while a load.data hook was still running")
	case <-time.After(100 * time.Millisecond):
	}

	close(hookRelease)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Set did not complete after the hook was released")
	}

	is.Eq(1, c.Int("first"))
	is.Eq(2, c.Int("second"))
}

func TestLoadDataInvalidatesCaches(t *testing.T) {
	c := New("load-cache", EnableCache)
	assert.NoErr(t, c.LoadData(map[string]any{
		"name":  "old",
		"names": []any{"old"},
	}))
	assert.Eq(t, "old", c.String("name"))
	assert.Eq(t, []string{"old"}, c.Strings("names"))

	assert.NoErr(t, c.LoadData(map[string]any{
		"name":  "new",
		"names": []any{"new"},
	}))
	assert.Eq(t, "new", c.String("name"))
	assert.Eq(t, []string{"new"}, c.Strings("names"))
}

func TestLoadStringsInvalidatesCaches(t *testing.T) {
	c := New("load-strings", EnableCache)
	assert.NoErr(t, c.LoadStrings(JSON, `{"name":"old"}`))
	assert.Eq(t, "old", c.String("name"))

	assert.NoErr(t, c.LoadStrings(JSON, `{"name":"new"}`))
	assert.Eq(t, "new", c.String("name"))
}

func TestClearDataInvalidatesCaches(t *testing.T) {
	c := New("clear-cache", EnableCache)
	assert.NoErr(t, c.LoadData(map[string]any{"name": "old"}))
	assert.Eq(t, "old", c.String("name"))

	c.ClearData()
	assert.Eq(t, "", c.String("name"))
}

// run with -race: the read caches are guarded by cacheMu, and a value read before a
// data change can not be stored into the cache after it.
func TestCacheConcurrentReadWrite(t *testing.T) {
	c := New("cache-race", EnableCache)
	assert.NoErr(t, c.LoadData(map[string]any{"name": "v0"}))

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			_ = c.String("name")
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			_ = c.Set("name", "v")
		}
	}()
	wg.Wait()

	assert.Eq(t, "v", c.String("name"))
}

// run with -race: the load paths change the data under the data lock.
func TestLoadDataConcurrentRead(t *testing.T) {
	c := New("load-race", EnableCache)
	assert.NoErr(t, c.LoadData(map[string]any{"name": "v0"}))

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = c.String("name")
			_ = c.Get("name")
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = c.LoadData(map[string]any{"name": "v1"})
		}
	}()
	wg.Wait()
}
