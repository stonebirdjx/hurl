// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: flag.go
// @Date: 2021/11/15 22:04
// @Desc: hurl parameter analysis
package configs

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var (
	Text      = flag.String("url", "", "enter a net/url text or local path")
	Walk      = flag.Bool("walk", false, "whether walk to the path")
	Mode      = flag.String("type", "", "walk the path type, value with file or dir")
	Regexp    = flag.String("re", "", "enter the regexp rule to match")
	Upload    = flag.String("upload", "", "upload local path to ftp")
	Download  = flag.String("download", "", "download ftp to local path")
	ReadBytes = flag.Int64("byte", maxBytes, "byte array max length default 10M")
	User      = flag.String("user", "", "enter the user name")
	Password  = flag.String("password", "", "enter the user password")
	Currency  = flag.Int("u", 5, "enter the number of concurrent")
	h         = flag.Bool("h", false, "print hurl tool help text")
	help      = flag.Bool("help", false, "print hurl tool help text")
	v         = flag.Bool("v", false, "print hurl tool version")
	version   = flag.Bool("version", false, "print hurl tool version")
)

func init() {
	flag.Parse()

	// 查看帮助信息, 正常退出
	if *h || *help {
		flag.Usage()
		os.Exit(0)
	}

	// 查看版本信息, 正常退出
	if *v || *version {
		fmt.Printf("%s %s for %s\n", toolName, toolVersion, runtime.GOOS)
		fmt.Println(versionInfo)
		os.Exit(0)
	}

}
