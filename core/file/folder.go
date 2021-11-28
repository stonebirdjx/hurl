// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: folder.go
// @Date: 2021/11/20 12:16
// @Desc:  deal file protocol path is folder
package file

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// 文件协议文件夹消息入口
func (b *BasicFile) folder() {
	switch b.Walk {
	case true:
		b.walkDir()
	default:
		b.readDir()
	}
}

// 读取文件夹信息
func (b *BasicFile) readDir() {
	fileInfos, err := ioutil.ReadDir(b.Path)
	if err != nil {
		log.Fatal(err)
	}
	// 先判断reg是否为nil
	switch b.Reg {
	case nil:
		b.readDirDealMode(fileInfos)
	default:
		b.readDirDealModeReg(fileInfos)
	}
}

// 没有正则直接执行
func (b *BasicFile) readDirDealMode(fileInfos []fs.FileInfo) {
	for _, fileInfo := range fileInfos {
		b.readDirPrint(fileInfo)
	}
}

// 有正则先执行正则
func (b *BasicFile) readDirDealModeReg(fileInfos []fs.FileInfo) {
	for _, fileInfo := range fileInfos {
		if b.Reg.FindString(fileInfo.Name()) == "" {
			continue
		}
		b.readDirPrint(fileInfo)
	}
}

// 基础打印方法
func (b *BasicFile) readDirPrint(fileInfo fs.FileInfo) {
	fType := ""
	switch b.Mode {
	case "file":
		if fileInfo.IsDir() {
			return
		}
	case "dir":
		if !fileInfo.IsDir() {
			return
		}
		fType = "<DIR>"
	default:
		if fileInfo.IsDir() {
			fType = "<DIR>"
		}
	}
	fmt.Printf("%s %5s %15d %s\n",
		fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		fType,
		fileInfo.Size(),
		fileInfo.Name(),
	)
}

// walk文件夹信息
func (b *BasicFile) walkDir() {
	err := filepath.Walk(b.Path, b.visit)
	if err != nil {
		log.Fatal(err)
	}
}

// walk函数的visit子函数
func (b *BasicFile) visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	switch b.Mode {
	case "file":
		if info.IsDir() {
			return nil
		}
	case "dir":
		if !info.IsDir() {
			return nil
		}
	}
	if b.Reg != nil {
		if b.Reg.FindString(path) != "" {
			fmt.Println(path)
		}
	} else {
		fmt.Println(path)
	}
	return nil
}
