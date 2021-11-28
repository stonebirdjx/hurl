// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: config.go
// @Date: 2021/11/15 22:00
// @Desc: hurl some static config
package configs

const (
	toolName    = "everything hurl"
	toolVersion = "v2.0.0"
	maxBytes    = 1 << 20
	File        = "file"
	FileDir     = "<DIR>"
	Dir         = "dir"
	Folder      = "folder"
	EmptyString = ""
	TimeFormat  = "2006-01-02 15:04:05"
	FileScheme  = File
	FtpScheme   = "ftp"
	SftpScheme  = "sftp"
	HttpScheme  = "http"
	HttpsScheme = "https"
)

var versionInfo = `
Protocols:file, ftp, sftp, http, https;
EMail:1245863260@qq.com, g1245863260@gmail.com;
Github:https://github.com/stonebirdjx;
Gitee:https://gitee.com/stonebirdjx;`
