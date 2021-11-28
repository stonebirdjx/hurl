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
	"io/fs"
	"log"
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
			if bsf.Reg.FindString(bs.Text()) != "" {
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
	switch bsf.Reg {
	case nil:
		bsf.readDirDealMode(fileInfos)
	default:
		bsf.readDirDealModeReg(fileInfos)
	}

}

func (bsf *BasicSftp) readDirDealMode(fileInfos []fs.FileInfo) {
	for _, fileInfo := range fileInfos {
		bsf.readDirPrint(fileInfo)
	}

}
func (bsf *BasicSftp) readDirDealModeReg(fileInfos []fs.FileInfo) {
	for _, fileInfo := range fileInfos {
		if bsf.Reg.FindString(fileInfo.Name()) == "" {
			continue
		}
		bsf.readDirPrint(fileInfo)
	}
}

func (bsf *BasicSftp) readDirPrint(fileInfo fs.FileInfo) {
	fType := ""
	switch bsf.Mode {
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

func (bsf *BasicSftp) walkDir() {

}
