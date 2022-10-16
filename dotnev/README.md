# Dotenv

Package `dotenv` that supports importing data from files (eg `.env`) to ENV

> Deprecated: please use package https://github.com/gookit/ini/v2/dotenv

## Install

```shell
go get github.com/gookit/ini/v2/dotenv
```

### Usage

```go
package main

import (
	"fmt"

	"github.com/gookit/ini/v2/dotenv"
)

func main() {
	err := dotenv.Load("./", ".env")
	if err != nil {
        fmt.Println(err)
	}
}
```

### Read Env

```go
val := dotenv.Get("ENV_KEY")
// Or use 
// val := os.Getenv("ENV_KEY")

// get int value
intVal := dotenv.Int("LOG_LEVEL")

// get bool value
blVal := dotenv.Bool("OPEN_DEBUG")

// with default value
val := dotenv.Get("ENV_KEY", "default value")
```
