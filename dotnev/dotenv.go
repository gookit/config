// Package dotnev provide load .env data to os ENV
//
// Deprecated: please use github.com/gookit/ini/v2/dotenv
package dotnev

import "github.com/gookit/ini/v2/dotenv"

// Load and parse .env file data to os ENV.
//
// Usage:
//
//	dotenv.Load("./", ".env")
//
// Deprecated: please use github.com/gookit/ini/v2/dotenv
//
//goland:noinspection ALL
var (
	Load        = dotenv.Load
	LoadExists  = dotenv.LoadExists
	LoadFromMap = dotenv.LoadFromMap

	LoadedData      = dotenv.LoadedData
	ClearLoaded     = dotenv.ClearLoaded
	DontUpperEnvKey = dotenv.DontUpperEnvKey
)

// Get ENV value by key
//
// Deprecated: please use github.com/gookit/ini/v2/dotenv
//
//goland:noinspection ALL
var (
	Get  = dotenv.Get
	Bool = dotenv.Bool
	Int  = dotenv.Int
)
