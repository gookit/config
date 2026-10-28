package config

import (
	"sync"
	"testing"
	"time"
)

func TestInspectionReadsConcurrentSetData(t *testing.T) {
	for _, tc := range []struct {
		name string
		read func(*Config) bool
	}{
		{"Exists", func(c *Config) bool { return c.Exists("profile.name") }},
		{"Keys", func(c *Config) bool { return len(c.Keys()) == 1 }},
		{"IsEmpty", func(c *Config) bool { return !c.IsEmpty() }},
		{"GetValue control", func(c *Config) bool {
			_, ok := c.GetValue("profile.name")
			return ok
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New("inspection-race")
			c.SetData(map[string]any{"profile": map[string]any{"name": "initial"}})
			start := make(chan struct{})
			failed := make(chan struct{}, 1)
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				<-start
				for i := 0; i < 200; i++ {
					if !tc.read(c) {
						failed <- struct{}{}
						return
					}
				}
			}()
			go func() {
				defer wg.Done()
				<-start
				for i := 0; i < 200; i++ {
					c.SetData(map[string]any{"profile": map[string]any{"name": "updated"}})
				}
			}()
			close(start)
			wg.Wait()
			select {
			case <-failed:
				t.Fatal("inspection did not observe a complete configuration")
			default:
			}
			if c.String("profile.name") != "updated" {
				t.Fatal("final configuration was not retained")
			}
		})
	}
}

func TestInspectionReadsInHooks(t *testing.T) {
	for _, tc := range []struct {
		event string
		write func(*Config) error
	}{
		{OnSetValue, func(c *Config) error { return c.Set("profile.name", "updated") }},
		{OnSetData, func(c *Config) error {
			c.SetData(map[string]any{"profile": map[string]any{"name": "updated"}})
			return nil
		}},
		{OnLoadData, func(c *Config) error {
			return c.LoadData(map[string]any{"profile": map[string]any{"name": "updated"}})
		}},
		{OnCleanData, func(c *Config) error { c.ClearData(); return nil }},
	} {
		t.Run(tc.event, func(t *testing.T) {
			observed := make(chan [3]bool, 1)
			c := New("inspection-hook", WithHookFunc(func(event string, cfg *Config) {
				if event == tc.event {
					observed <- [3]bool{cfg.Exists("profile.name"), len(cfg.Keys()) == 1, cfg.IsEmpty()}
				}
			}))
			done := make(chan error, 1)
			go func() { done <- tc.write(c) }()
			select {
			case err := <-done:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("write deadlocked while its hook inspected the configuration")
			}
			want := [3]bool{true, true, false}
			if tc.event == OnCleanData {
				want = [3]bool{false, false, true}
			}
			select {
			case got := <-observed:
				if got != want {
					t.Fatalf("hook inspection = %v, want %v", got, want)
				}
			default:
				t.Fatal("write did not invoke its hook")
			}
		})
	}
}
