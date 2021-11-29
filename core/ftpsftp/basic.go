// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: basic.go
// @Date: 2021/11/21 17:30
// @Desc: sftp or ftp protocol basic information
package ftpsftp

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
)

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

type transport struct {
	name     string // 名称
	tp       string // dir, file, link会被当file传上传
	size     uint64 // 文件大小
	relative string // 相对路径
	site     string // 本地位置
}

var trChan = make(chan transport)
var mutex sync.Mutex

// sftp,ftp 下载时本地必须是文件夹,不存在时创建
func downloadPathIsDir(local string) {
	localInfo, err := os.Stat(local)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(local, 0644)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	} else if !localInfo.IsDir() {
		log.Fatal("local path is not a dir")
	}
}

// 检查或者创建本地文件夹
func cmLocalDir(dir string) error {
	mutex.Lock()
	defer mutex.Unlock()
	fileInfo, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0644)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !fileInfo.IsDir() {
		return fmt.Errorf("%s exist, but is not a dir", dir)
	}
	return nil
}
