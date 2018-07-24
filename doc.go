/*
Package config is a go config management implement. support YAML,TOML,JSON,INI,HCL format.

Source code and other details for the project are available at GitHub:

	https://github.com/gookit/config

JSON format content example:

	{
		"name": "app",
		"debug": false,
		"baseKey": "value",
		"age": 123,
		"envKey": "${SHELL}",
		"envKey1": "${NotExist|defValue}",
		"map1": {
			"key": "val",
			"key1": "val1",
			"key2": "val2"
		},
		"arr1": [
			"val",
			"val1",
			"val2"
		],
		"lang": {
			"dir": "res/lang",
			"defLang": "en",
			"allowed": {
				"en": "val",
				"zh-CN": "val2"
			}
		}
	}

Usage please see example(more example please see examples folder in the lib):

*/
package config
