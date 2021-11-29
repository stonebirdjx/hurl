// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: folder.go
// @Date: 2021/11/20 12:16
// @Desc:  deal file protocol path is folder
package file

import (
	"fmt"
	"hurl/configs"
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
		b.walkDir() // walk文件夹目录
	default:
		b.readDir() // 读取当前文件夹
	}
}

// 读取当前文件夹信息
func (b *BasicFile) readDir() {
	fileInfos, err := ioutil.ReadDir(b.Path)
	if err != nil {
		log.Fatal(err)
	}

	// 先判断reg是否为nil,减少判断次数
	switch b.Reg {
	case nil:
		b.readDirDealMode(fileInfos) // 没配正则执行
	default:
		b.readDirDealModeReg(fileInfos) // 正则执行
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
		if b.Reg.FindString(fileInfo.Name()) == configs.EmptyString {
			continue
		}
		b.readDirPrint(fileInfo)
	}
}

// 基础打印方法
func (b *BasicFile) readDirPrint(fileInfo fs.FileInfo) {
	fType := ""
	switch b.Mode {
	case configs.File:
		if fileInfo.IsDir() {
			return
		}
	case configs.Dir:
		if !fileInfo.IsDir() {
			return
		}
		fType = configs.FileDir
	default:
		if fileInfo.IsDir() {
			fType = configs.FileDir
		}
	}
	fmt.Printf("%s %5s %15d %s\n",
		fileInfo.ModTime().Format(configs.TimeFormat),
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
	case configs.File:
		if info.IsDir() {
			return nil
		}
	case configs.Dir:
		if !info.IsDir() {
			return nil
		}
	}
	if b.Reg != nil {
		if b.Reg.FindString(path) != configs.EmptyString {
			fmt.Println(path)
		}
	} else {
		fmt.Println(path)
	}
	return nil
}
