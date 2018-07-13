# ini config read

```go
package main

import (
	newini "github.com/ochinchina/go-ini"
	kini "github.com/knq/ini"
    dini "github.com/dombenson/go-ini"
)

    // "gopkg.in/ini.v1" 多个数据源合并后，获取section数据有误
    
	Cfg, err := ini.Load("conf/app.ini", envFile)
	// 只读时提升性能
	Cfg.BlockMode = false

	// newini "github.com/ochinchina/go-ini"
	config := newini.Load("conf/app.ini", envFile)
	rds, err := config.GetSection("redis")
	fmt.Printf("merged: %s \n rds server: %s \n", config.String(), rds.Key("server").String())

	// kini "github.com/knq/ini" 只能单个文件
	cfg1, err := kini.LoadFile("conf/app.ini")
	fmt.Printf("merged: %s \n rds server: %s \n", cfg1.String(), cfg1.GetSection("redis").Get("server"))

	// dini "github.com/dombenson/go-ini"
	cfg2, err := dini.LoadFile("conf/app.ini")
	cfg3, err := dini.LoadFile(envFile)

	cfg2.WriteTo(os.Stdout)

	fmt.Printf("rds server: %+v \n", cfg2.Values("redis"))

	// cfg3 copy to cfg2 通过copy，可以合并多个文件
	cfg3.Copy(cfg2)

	cfg2.WriteTo(os.Stdout)
	fmt.Printf("log conf: %+v \n", cfg2.Values("log"))

```
