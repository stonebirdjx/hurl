// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: hurl.go
// @Date: 2021/11/14 20:42
// @Desc: hurl tool v2.0 main function
package main

import (
	"hurl/configs"
	"hurl/core"
	"hurl/core/file"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {
	text := strings.TrimSpace(*configs.Text)
	if text == "" {
		return
	}

	u, err := url.Parse(text)
	if err != nil {
		log.Fatal(err)
	}

	switch u.Scheme {
	case configs.FileScheme, configs.EmptyString:
		core.FileHandle(u)
	case configs.HttpScheme, configs.HttpsScheme:
		core.HttpHandle(u)
	case configs.FtpScheme, configs.SftpScheme:
		core.FtpSftpHandle(u)
	default:
		_, err := os.Stat(text)
		if err != nil {
			log.Fatal("hurl current support file, http, https, ftp, sftp protocol")
		}

		basicFile, err := file.NewBasicFiler(text)
		if err != nil {
			log.Fatal(err)
		}

		basicFile.Entrance()
	}
}
