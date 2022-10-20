//go:build !go1.18
// +build !go1.18

package config

// alias of interface{}, use for go < 1.18
type any = interface{}
