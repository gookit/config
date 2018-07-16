/*
golang application config manage implement. support yaml,toml,json format.

Source code and other details for the project are available at GitHub:

- https://github.com/gookit/config

Yaml format content:

	name: app2
	debug: false
	baseKey: value2

	map1:
	    key: val2
	    key2: val20

	arr1:
	    - val1
	    - val21

Toml format content:

	title = "TOML Example"
	name = "app"

	envKey = "${SHELL}"
	envKey1 = "${NotExist|defValue}"

	arr1 = [
	  "alpha",
	  "omega"
	]

	[map1]
	name = "inhere"
	org = "GitHub"

Usage please see examples:

- yaml format example and usage, see [yaml-example](/github.com/gookit/config/yaml#example-package)
- toml format example and usage, see [toml-example](/github.com/gookit/config/toml#example-package)

 */
package config
