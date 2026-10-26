package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
)

const stringsEnvToken = "${CONFIG_STRINGS_PRESENT}"

type stringsEnvStringer struct{}

func (stringsEnvStringer) String() string { return stringsEnvToken }

func stringsEnvChange(t *testing.T, name string, value *string) {
	t.Helper()
	var err error
	if value == nil {
		err = os.Unsetenv(name)
	} else {
		err = os.Setenv(name, *value)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func stringsEnvSetup(t *testing.T) {
	t.Helper()
	for _, name := range []string{"CONFIG_STRINGS_PRESENT", "CONFIG_STRINGS_MISSING", "CONFIG_STRINGS_EMPTY"} {
		key := name
		value, present := os.LookupEnv(key)
		t.Cleanup(func() {
			if present {
				stringsEnvChange(t, key, &value)
				return
			}
			stringsEnvChange(t, key, nil)
		})
	}
	value, empty := "resolved", ""
	stringsEnvChange(t, "CONFIG_STRINGS_PRESENT", &value)
	stringsEnvChange(t, "CONFIG_STRINGS_MISSING", nil)
	stringsEnvChange(t, "CONFIG_STRINGS_EMPTY", &empty)
}

func stringsEnvEqual(t *testing.T, label string, want, got any) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%s: want %#v, got %#v", label, want, got)
	}
}

func stringsEnvInput(kind string, values []string) any {
	if kind == "strings" {
		return values
	}
	var result []any
	for _, value := range values {
		result = append(result, value)
	}
	return result
}

func stringsEnvConfig(t *testing.T, value any, parse, cache bool) *Config {
	t.Helper()
	c := New("strings-env", func(opts *Options) {
		opts.ParseEnv = parse
		opts.EnableCache = cache
	})
	if err := c.Set("list", value); err != nil {
		t.Fatal(err)
	}
	return c
}

func stringsEnvRaw(t *testing.T, c *Config, input any) {
	t.Helper()
	value, ok := c.GetValue("list")
	stringsEnvEqual(t, "GetValue found", true, ok)
	stringsEnvEqual(t, "GetValue raw", input, value)
	stringsEnvEqual(t, "Get raw", input, c.Get("list"))
	stringsEnvEqual(t, "Data raw", input, c.Data()["list"])
}

func TestStringsParseEnvValues(t *testing.T) {
	stringsEnvSetup(t)
	raw := []string{
		stringsEnvToken,
		"${CONFIG_STRINGS_MISSING|defaulted}",
		"${ CONFIG_STRINGS_EMPTY | empty-default }",
		"${CONFIG_STRINGS_MISSING}",
		"${CONFIG_STRINGS_EMPTY}",
		"pre/${CONFIG_STRINGS_PRESENT}/${ CONFIG_STRINGS_MISSING | tail }",
		"literal",
		"",
	}
	want := []string{"resolved", "defaulted", "empty-default", "", "", "pre/resolved/tail", "literal", ""}
	for _, kind := range []string{"strings", "any"} {
		t.Run(kind, func(t *testing.T) {
			c := New("strings-env-values", ParseEnv)
			if err := c.Set("nested.list", stringsEnvInput(kind, raw)); err != nil {
				t.Fatal(err)
			}
			stringsEnvEqual(t, "expanded values", want, c.Strings("nested.list"))
			stringsEnvEqual(t, "getter error", nil, c.Error())
		})
	}
}

func TestStringsParseEnvDisabled(t *testing.T) {
	stringsEnvSetup(t)
	for _, kind := range []string{"strings", "any"} {
		for _, cache := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/cache_%t", kind, cache), func(t *testing.T) {
				raw := []string{stringsEnvToken, "${CONFIG_STRINGS_MISSING|defaulted}", "literal", ""}
				input := stringsEnvInput(kind, raw)
				c := stringsEnvConfig(t, input, false, cache)
				stringsEnvEqual(t, "first raw read", raw, c.Strings("list"))
				value := "changed"
				stringsEnvChange(t, "CONFIG_STRINGS_PRESENT", &value)
				result := c.Strings("list")
				stringsEnvEqual(t, "second raw read", raw, result)
				result[0] = "mutated"
				if kind == "strings" {
					stringsEnvEqual(t, "disabled alias", "mutated", raw[0])
					stringsEnvEqual(t, "stored alias", "mutated", c.Get("list").([]string)[0])
					return
				}
				stringsEnvEqual(t, "any caller unchanged", stringsEnvToken, input.([]any)[0])
				stringsEnvRaw(t, c, []any{stringsEnvToken, "${CONFIG_STRINGS_MISSING|defaulted}", "literal", ""})
			})
		}
	}
	t.Run("any_conversions", func(t *testing.T) {
		for _, cache := range []bool{false, true} {
			t.Run(fmt.Sprintf("cache_%t", cache), func(t *testing.T) {
				c := stringsEnvConfig(t, []any{stringsEnvToken, 12, true, nil}, false, cache)
				want := []string{stringsEnvToken, "12", "true", ""}
				stringsEnvEqual(t, "first disabled conversion", want, c.Strings("list"))
				stringsEnvEqual(t, "second disabled conversion", want, c.Strings("list"))
			})
		}
	})
}

func TestStringsParseEnvShapes(t *testing.T) {
	stringsEnvSetup(t)
	cases := []struct {
		name      string
		key       string
		value     any
		want      []string
		wantError bool
	}{
		{"missing", "absent", nil, nil, false},
		{"empty_key", "", nil, nil, true},
		{"nil_strings", "list", []string(nil), nil, false},
		{"empty_strings", "list", []string{}, []string{}, false},
		{"nil_any", "list", []any(nil), nil, false},
		{"empty_any", "list", []any{}, nil, false},
		{"untyped_nil", "list", nil, nil, true},
		{"scalar", "list", "literal", nil, true},
		{"scalar_token", "list", stringsEnvToken, nil, true},
		{"ints", "list", []int{1}, nil, true},
		{"map", "list", map[string]any{"x": 1}, nil, true},
		{"array", "list", [1]string{"x"}, nil, true},
	}
	for _, parse := range []bool{false, true} {
		for _, tc := range cases {
			t.Run(fmt.Sprintf("parse_%t/%s", parse, tc.name), func(t *testing.T) {
				c := stringsEnvConfig(t, tc.value, parse, false)
				stringsEnvEqual(t, "slice shape", tc.want, c.Strings(tc.key))
				err := c.Error()
				stringsEnvEqual(t, "error presence", tc.wantError, err != nil)
				if tc.key == "" {
					stringsEnvEqual(t, "empty key error", true, errors.Is(err, ErrKeyIsEmpty))
				}
			})
		}
	}
}

func stringsEnvPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("unsupported element conversion did not panic")
		}
	}()
	fn()
}

func TestStringsParseEnvConversions(t *testing.T) {
	stringsEnvSetup(t)
	t.Run("non_string", func(t *testing.T) {
		for _, parse := range []bool{false, true} {
			t.Run(fmt.Sprintf("parse_%t", parse), func(t *testing.T) {
				token := stringsEnvToken
				input := []any{12, true, 1.25, nil, []byte(token), &token, stringsEnvStringer{}, errors.New(token), []byte(nil)}
				c := stringsEnvConfig(t, input, parse, false)
				stringsEnvEqual(t, "non-string conversions", []string{"12", "true", "1.25", "", token, token, token, token, ""}, c.Strings("list"))
			})
		}
	})
	t.Run("mixed", func(t *testing.T) {
		c := stringsEnvConfig(t, []any{stringsEnvToken, 12, true, nil}, true, false)
		stringsEnvEqual(t, "mixed conversions", []string{"resolved", "12", "true", ""}, c.Strings("list"))
	})
	t.Run("unsupported", func(t *testing.T) {
		for _, parse := range []bool{false, true} {
			t.Run(fmt.Sprintf("parse_%t", parse), func(t *testing.T) {
				c := stringsEnvConfig(t, []any{struct{}{}}, parse, false)
				stringsEnvPanic(t, func() { c.Strings("list") })
			})
		}
	})
}

func TestStringsParseEnvOwnership(t *testing.T) {
	stringsEnvSetup(t)
	for _, kind := range []string{"strings", "any"} {
		for _, cache := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/cache_%t", kind, cache), func(t *testing.T) {
				input := stringsEnvInput(kind, []string{stringsEnvToken, "literal"})
				c := stringsEnvConfig(t, input, true, cache)
				result := c.Strings("list")
				stringsEnvEqual(t, "resolved result", []string{"resolved", "literal"}, result)
				wantRaw := stringsEnvInput(kind, []string{stringsEnvToken, "literal"})
				stringsEnvEqual(t, "caller before mutation", wantRaw, input)
				stringsEnvRaw(t, c, wantRaw)
				result[0] = "mutated"
				stringsEnvEqual(t, "caller after mutation", wantRaw, input)
				stringsEnvRaw(t, c, wantRaw)
			})
		}
	}
}

func TestStringsParseEnvFreshEnvironment(t *testing.T) {
	stringsEnvSetup(t)
	for _, kind := range []string{"strings", "any"} {
		t.Run(kind, func(t *testing.T) {
			token := "${CONFIG_STRINGS_PRESENT|defaulted}"
			input := stringsEnvInput(kind, []string{token})
			c := stringsEnvConfig(t, input, true, false)
			value := "first"
			stringsEnvChange(t, "CONFIG_STRINGS_PRESENT", &value)
			first := c.Strings("list")
			stringsEnvEqual(t, "first environment", []string{"first"}, first)
			value = "second"
			stringsEnvChange(t, "CONFIG_STRINGS_PRESENT", &value)
			second := c.Strings("list")
			stringsEnvEqual(t, "second environment", []string{"second"}, second)
			stringsEnvEqual(t, "prior result retained", []string{"first"}, first)
			value = ""
			stringsEnvChange(t, "CONFIG_STRINGS_PRESENT", &value)
			stringsEnvEqual(t, "empty selects default", []string{"defaulted"}, c.Strings("list"))
			stringsEnvChange(t, "CONFIG_STRINGS_PRESENT", nil)
			stringsEnvEqual(t, "absent selects default", []string{"defaulted"}, c.Strings("list"))
			second[0] = "mutated"
			stringsEnvEqual(t, "fresh after result mutation", []string{"defaulted"}, c.Strings("list"))
			stringsEnvEqual(t, "caller raw", stringsEnvInput(kind, []string{token}), input)
			stringsEnvRaw(t, c, stringsEnvInput(kind, []string{token}))
		})
	}
}

func stringsEnvCachedRead(t *testing.T, c *Config, want string) {
	t.Helper()
	stringsEnvEqual(t, "scalar cache", want, c.String("scalar"))
	stringsEnvEqual(t, "list cache", []string{want}, c.Strings("list"))
}

// TestStringsParseEnvCache pins the cache semantics: a cache hit returns the value that
// was expanded on the first read, so an environment change is only picked up after
// ClearCaches() or a new Set() - the same as the String() cache.
func TestStringsParseEnvCache(t *testing.T) {
	stringsEnvSetup(t)
	for _, kind := range []string{"strings", "any"} {
		t.Run(kind, func(t *testing.T) {
			input := stringsEnvInput(kind, []string{stringsEnvToken})
			c := stringsEnvConfig(t, input, true, true)
			if err := c.Set("scalar", stringsEnvToken); err != nil {
				t.Fatal(err)
			}
			value := "first"
			stringsEnvChange(t, "CONFIG_STRINGS_PRESENT", &value)
			stringsEnvCachedRead(t, c, "first")
			value = "second"
			stringsEnvChange(t, "CONFIG_STRINGS_PRESENT", &value)
			stringsEnvCachedRead(t, c, "first")
			c.ClearCaches()
			stringsEnvCachedRead(t, c, "second")
			stringsEnvEqual(t, "scalar raw", stringsEnvToken, c.Get("scalar"))
			stringsEnvRaw(t, c, stringsEnvInput(kind, []string{stringsEnvToken}))
			value = "replacement"
			stringsEnvChange(t, "CONFIG_STRINGS_MISSING", &value)
			if err := c.Set("list", stringsEnvInput(kind, []string{"${CONFIG_STRINGS_MISSING}"})); err != nil {
				t.Fatal(err)
			}
			stringsEnvEqual(t, "Set invalidates", []string{"replacement"}, c.Strings("list"))
		})
	}
}

func stringsEnvDefault(t *testing.T) {
	t.Helper()
	previous := dc
	dc = New("strings-env-default")
	t.Cleanup(func() { dc = previous })
	WithOptions(ParseEnv)
}

func TestStringsParseEnvEntrypoints(t *testing.T) {
	stringsEnvSetup(t)
	t.Run("package_json", func(t *testing.T) {
		stringsEnvDefault(t)
		if err := LoadStrings(JSON, `{"nested":{"list":["${CONFIG_STRINGS_PRESENT}","${CONFIG_STRINGS_MISSING|defaulted}"]}}`); err != nil {
			t.Fatal(err)
		}
		want := []string{"resolved", "defaulted"}
		stringsEnvEqual(t, "package getter", want, Strings("nested.list"))
		stringsEnvEqual(t, "default getter", want, Default().Strings("nested.list"))
		stringsEnvEqual(t, "package raw getter", []any{stringsEnvToken, "${CONFIG_STRINGS_MISSING|defaulted}"}, Get("nested.list"))
	})
	t.Run("constructors", func(t *testing.T) {
		cases := []struct {
			name   string
			config *Config
		}{
			{"NewWithOptions", NewWithOptions("strings-env-alias", ParseEnv)},
			{"NewGeneric", NewGeneric("strings-env-generic")},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if err := tc.config.Set("list", []string{stringsEnvToken}); err != nil {
					t.Fatal(err)
				}
				stringsEnvEqual(t, "constructor getter", []string{"resolved"}, tc.config.Strings("list"))
			})
		}
	})
	t.Run("scalar_controls", func(t *testing.T) {
		stringsEnvDefault(t)
		value := "one,two"
		stringsEnvChange(t, "CONFIG_STRINGS_PRESENT", &value)
		if err := Set("scalar", stringsEnvToken); err != nil {
			t.Fatal(err)
		}
		stringsEnvEqual(t, "instance split", []string{"one", "two"}, Default().StringsBySplit("scalar", ","))
		stringsEnvEqual(t, "package split", []string{"one", "two"}, StringsBySplit("scalar", ","))
		stringsEnvEqual(t, "scalar String", "one,two", String("scalar"))
		stringsEnvEqual(t, "scalar MustString", "one,two", MustString("scalar"))
		stringsEnvEqual(t, "literal missing default", stringsEnvToken, String("absent", stringsEnvToken))
	})
}

func stringsEnvReadWorker(c *Config, start <-chan struct{}, failures chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()
	<-start
	want := []string{"resolved", "literal"}
	for i := 0; i < 50; i++ {
		got := c.Strings("list")
		if !reflect.DeepEqual(want, got) {
			failures <- got
			return
		}
	}
}

func stringsEnvConcurrentRead(t *testing.T, kind string, readonly, cache bool) {
	t.Helper()
	input := stringsEnvInput(kind, []string{stringsEnvToken, "literal"})
	c := stringsEnvConfig(t, input, true, cache)
	if readonly {
		c.Readonly()
	}
	start := make(chan struct{})
	failures := make(chan []string, 8)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go stringsEnvReadWorker(c, start, failures, &wg)
	}
	close(start)
	wg.Wait()
	close(failures)
	for got := range failures {
		stringsEnvEqual(t, "concurrent resolved read", []string{"resolved", "literal"}, got)
	}
	stringsEnvEqual(t, "warm resolved read", []string{"resolved", "literal"}, c.Strings("list"))
	stringsEnvEqual(t, "concurrent caller raw", stringsEnvInput(kind, []string{stringsEnvToken, "literal"}), input)
	stringsEnvRaw(t, c, stringsEnvInput(kind, []string{stringsEnvToken, "literal"}))
}

func TestStringsParseEnvConcurrentReaders(t *testing.T) {
	stringsEnvSetup(t)
	for _, kind := range []string{"strings", "any"} {
		for _, readonly := range []bool{false, true} {
			for _, cache := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/readonly_%t/cache_%t", kind, readonly, cache), func(t *testing.T) {
					stringsEnvConcurrentRead(t, kind, readonly, cache)
				})
			}
		}
	}
}
