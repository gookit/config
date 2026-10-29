package config_test

import (
	"github.com/gookit/config/v2"
	"testing"
	"time"
)

func TestStructureNegativeDurations(t *testing.T) {
	for _, value := range []string{"-1ns", "-1us", "-1µs", "-1ms", "-1s", "-1m", "-1h", "-1h30m", "-0s", "0s", "1h30m", "100ms"} {
		t.Run(value, func(t *testing.T) {
			c := config.New("duration", config.ParseTime)
			c.SetData(map[string]any{"timeout": value, "label": value})
			var dst struct {
				Timeout time.Duration `mapstructure:"timeout"`
				Label   string        `mapstructure:"label"`
			}
			if err := c.Decode(&dst); err != nil {
				t.Fatal(err)
			}
			want, err := time.ParseDuration(value)
			if err != nil {
				t.Fatal(err)
			}
			if dst.Timeout != want {
				t.Fatalf("duration: got %v, want %v", dst.Timeout, want)
			}
			if dst.Label != value {
				t.Fatalf("string field changed: got %q", dst.Label)
			}
		})
	}
}

func TestStructureSignedDurationFromEnvironment(t *testing.T) {
	t.Setenv("GOTOOL_CONFIG_SIGNED_DURATION", "-500ms")
	c := config.New("duration", config.ParseTime, config.ParseEnv)
	if err := c.LoadStrings(config.JSON, `{"timeout": "${GOTOOL_CONFIG_SIGNED_DURATION}"}`); err != nil {
		t.Fatal(err)
	}
	var dst struct {
		Timeout time.Duration `mapstructure:"timeout"`
	}
	if err := c.Decode(&dst); err != nil {
		t.Fatal(err)
	}
	if dst.Timeout != -500*time.Millisecond {
		t.Fatalf("duration: got %v", dst.Timeout)
	}
}

func TestStructureNegativeDurationNumericControls(t *testing.T) {
	for _, tc := range []struct {
		value     string
		want      int64
		wantError bool
	}{
		{"-42", -42, false}, {"+42", 42, false}, {"42", 0, true},
	} {
		t.Run(tc.value, func(t *testing.T) {
			c := config.New("numeric", config.ParseTime)
			c.SetData(map[string]any{"count": tc.value, "delay": tc.value})
			var dst struct {
				Count int64
				Delay time.Duration
			}
			err := c.Decode(&dst)
			if tc.wantError {
				// The existing positive-number duration path requires a unit.
				if err == nil {
					t.Fatal("expected the existing missing-unit error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if dst.Count != tc.want || int64(dst.Delay) != tc.want {
				t.Fatalf("got %+v, want %d", dst, tc.want)
			}
		})
	}
}
