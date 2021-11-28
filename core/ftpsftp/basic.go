// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: basic.go
// @Date: 2021/11/21 17:30
// @Desc: sftp or ftp protocol basic information
package ftpsftp

import "regexp"

type BasicApi interface {
	Read()
	Upload()
	Download()
}

// ftp、sftp基本信息结构体
type BasicStruct struct {
	Path     string // ftp、sftp路径，路径必须以/结尾
	Host     string // ftp、sftp host地址(hostname:port)
	User     string // ftp、sftp用户名
	Password string // ftp、sftp用户名、密码
	Walk     bool   // walk遍历ftp、sftp的path目录
	Mode     string // all file dir
	Reg      *regexp.Regexp
}

type Joint struct {
	BasicStruct
}

type BasicFtp Joint
type BasicSftp Joint
