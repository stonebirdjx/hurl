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
	Text    = flag.String("url", "", "enter a net/url text or local path")
	h       = flag.Bool("h", false, "print hurl tool help text")
	help    = flag.Bool("help", false, "print hurl tool help text")
	v       = flag.Bool("v", true, "print hurl tool version")
	version = flag.Bool("version", false, "print hurl tool version")
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
