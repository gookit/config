package main

import (
	"github.com/gookit/config"
	"github.com/gookit/config/yaml"
	"fmt"
)

func main() {
	config.SetOptions(&config.Options{})

	config.SetDecoder(config.Yaml, yaml.Decoder)

	config.LoadFiles("testdata/yml_base.yml")

	fmt.Printf("config data: \n %#v\n", config.Data())

	config.LoadFiles("testdata/yml_other.yml")
	// config.LoadFiles("testdata/yml_base.yml", "testdata/yml_other.yml")

	fmt.Printf("config data: \n %#v\n", config.Data())
}
