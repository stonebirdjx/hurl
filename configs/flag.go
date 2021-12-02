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
	Method    = flag.String("method", "GET", "enter http protocol request method")
	Ho        = flag.String("o", "", "enter file name to save http body")
	HI        = flag.Bool("I", false, "whether HEAD request http request")
	Hi        = flag.Bool("i", false, "whether print http request headers")
	HL        = flag.Bool("L", false, "whether to open the redirect request")
	Headers   = flag.String("headers", "", "set http request header,must json type")
	Data      = flag.String("data", "", "http request body data")
	MultiPart = flag.Bool("multipart", false, "whether send a multipart/form request")
	HFile     = flag.String("file", "", "input file name use to multipart/form")
	Crt       = flag.String("crt", "", "enter the crt file")
	Http2     = flag.Bool("http2", false, "whether use http 2 protocol")
	Http3     = flag.Bool("http3", false, "whether use http 3 protocol")
	Walk      = flag.Bool("walk", false, "whether walk to the path")
	Mode      = flag.String("type", "", "walk the path type, value with file or dir")
	Regexp    = flag.String("re", "", "enter the regexp rule to match")
	Upload    = flag.String("upload", "", "upload local path to ftp")
	Download  = flag.String("download", "", "download ftp to local path")
	ReadBytes = flag.Int64("byte", maxBytes, "byte array max length default 1M")
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
