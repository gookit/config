# config

[![GoDoc](https://godoc.org/github.com/gookit/config?status.svg)](https://godoc.org/github.com/gookit/config)

golang application config manage implement. 

- generic api `Get` `Set` ...
- support file format: `json`, `yaml`, `toml`
- support multi file load
- data override merge

## Godoc

- [godoc for gopkg](https://godoc.org/gopkg.in/gookit/config.v1)
- [godoc for github](https://godoc.org/github.com/gookit/config)

## Usage

> there are use `ini` config file

```text
conf/
    base.yml
    dev.yml
    test.yml
    ...
```

- init

```go
    import "github/gookit/config"

    languages := map[string]string{
        "en": "English",
        "zh-CN": "简体中文",
        "zh-TW": "繁体中文",
    }

    config.Init("conf/lang", "en", languages)
```

- usage

```go
    // translate from special language
    val := config.Tr("en", "key")

    // translate from default language
    val := config.DefTr("key")
```

## useful packages

### ini parse

- [go-ini/ini](https://github.com/go-ini/ini) ini parser and config manage
- [dombenson/go-ini](https://github.com/dombenson/go-ini) ini parser and config manage

### yaml

- [go-yaml](https://github.com/go-yaml/yaml) ini parser


### toml

- [go toml](https://github.com/BurntSushi/toml) toml parser

### data merge

- [mergo](github.com/imdario/mergo) merge data

## License

MIT
