// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: sftpread.go
// @Date: 2021/11/28 8:35
// @Desc: sftp read file or dir
package ftpsftp

import (
	"bufio"
	"fmt"
	walkfs "github.com/kr/fs"
	"hurl/configs"
	"io/fs"
	"log"
	"os"
	"strings"
)

// read 主函数处理，用于分流
func (bsf *BasicSftp) Read() {
	if strings.HasSuffix(bsf.Path, "/") {
		if bsf.Walk {
			bsf.walkDir()
		} else {
			bsf.readDir()
		}
	} else {
		bsf.readFile()
	}
}

// 阅读单个文件
func (bsf *BasicSftp) readFile() {
	c, err := bsf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	f, err := c.Open(bsf.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bs := bufio.NewScanner(f)
	switch bsf.Reg {
	case nil:
		for bs.Scan() {
			fmt.Println(bs.Text())
		}
	default:
		for bs.Scan() {
			if bsf.Reg.FindString(bs.Text()) != configs.EmptyString {
				fmt.Println(bs.Text())
			}
		}
	}
}

// 读取文件夹
func (bsf *BasicSftp) readDir() {
	c, err := bsf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	fileInfos, err := c.ReadDir(bsf.Path)
	if err != nil {
		log.Fatal(err)
	}

	// 先判断是否有正则
	switch bsf.Reg {
	case nil:
		bsf.readDirDealMode(fileInfos)
	default:
		bsf.readDirDealModeReg(fileInfos)
	}

}

// 没正则直接打印
func (bsf *BasicSftp) readDirDealMode(fileInfos []fs.FileInfo) {
	for _, fileInfo := range fileInfos {
		bsf.readDirPrint(fileInfo)
	}
}

// 判断正则信息后打印
func (bsf *BasicSftp) readDirDealModeReg(fileInfos []fs.FileInfo) {
	for _, fileInfo := range fileInfos {
		if bsf.Reg.FindString(fileInfo.Name()) == configs.EmptyString {
			continue
		}
		bsf.readDirPrint(fileInfo)
	}
}

// 基础打印
func (bsf *BasicSftp) readDirPrint(fileInfo fs.FileInfo) {
	fType := ""
	switch bsf.Mode {
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

// walk 遍历文件夹
func (bsf *BasicSftp) walkDir() {
	c, err := bsf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	walker := c.Walk(bsf.Path)
	switch bsf.Reg {
	case nil:
		bsf.walkDirMode(walker)
	default:
		bsf.walkDirModeReg(walker)
	}
}

// 没有正则的情况下遍历文件夹
func (bsf *BasicSftp) walkDirMode(walker *walkfs.Walker) {
	for walker.Step() {
		_, next := bsf.walkDirModeBase(walker)
		if !next {
			continue
		}
		fmt.Println(walker.Path())
	}
}

// 有正则的情况下遍历文件夹
func (bsf *BasicSftp) walkDirModeReg(walker *walkfs.Walker) {
	for walker.Step() {
		fileInfo, next := bsf.walkDirModeBase(walker)
		if !next {
			continue
		}
		if bsf.Reg.FindString(fileInfo.Name()) != configs.EmptyString {
			fmt.Println(walker.Path())
		}
	}
}

// walk遍历文件夹公共方法
func (bsf *BasicSftp) walkDirModeBase(walker *walkfs.Walker) (os.FileInfo, bool) {
	fileInfo := walker.Stat()
	switch bsf.Mode {
	case configs.File:
		if fileInfo.IsDir() {
			return fileInfo, false
		}
	case configs.Dir:
		if !fileInfo.IsDir() {
			return fileInfo, false
		}
	}
	return fileInfo, true
}
