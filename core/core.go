// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: core.go
// @Date: 2021/11/20 10:49
// @Desc: core shunt layer
package core

import (
	"hurl/configs"
	"hurl/core/file"
	"hurl/core/ftpsftp"
	"hurl/shares"
	"log"
	"net/url"
	"strings"
)

// 文件协议消息处理者
// 传入类型 *url.URL
func FileHandle(u *url.URL) {
	path := u.Path
	basicFile, err := file.NewBasicFiler(path)
	if err != nil {
		log.Fatal(err)
	}
	basicFile.Entrance()
}

// sftp和ftp消息处理者
// 传入类型 *url.URL
func FtpSftpHandle(u *url.URL) {
	userName := u.User.Username()
	if userName == configs.EmptyString {
		userName = *configs.User
		if userName == configs.EmptyString {
			log.Fatalf("can not get ftp user")
		}
	}

	passWord, ok := u.User.Password()
	if !ok {
		passWord = *configs.Password
	}

	path := u.Path
	// 以//开头表绝对路径，否则是相对路径
	if !strings.HasPrefix(path, "//") {
		path = "." + path
	}

	// 判断是否使用正则表达式
	reg, err := shares.IfReg()
	if err != nil {
		log.Fatal(err)
	}

	var api ftpsftp.BasicApi
	basicStruct := ftpsftp.BasicStruct{
		Path:     path,
		Host:     u.Host,
		User:     userName,
		Password: passWord,
		Walk:     *configs.Walk,
		Mode:     *configs.Mode,
		Reg:      reg,
	}

	switch u.Scheme {
	case configs.FtpScheme:
		api = &ftpsftp.BasicFtp{
			BasicStruct: basicStruct,
		}
	case configs.SftpScheme:
		api = &ftpsftp.BasicSftp{
			BasicStruct: basicStruct,
		}
	default:
		log.Fatal("scheme is not ftp or sftp")
	}

	switch {
	case strings.TrimSpace(*configs.Download) != "":
		api.Download()
	case strings.TrimSpace(*configs.Upload) != "":
		api.Upload()
	default:
		api.Read()
	}
}
