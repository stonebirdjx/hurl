// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: ftpread.go
// @Date: 2021/11/28 17:22
// @Desc: ftp read file or dir
package ftpsftp

import (
	"bufio"
	"fmt"
	"github.com/jlaffaye/ftp"
	"log"
	"strings"
)

func (bf *BasicFtp) Read() {
	if strings.HasSuffix(bf.Path, "/") {
		if bf.Walk {
			bf.walkDir()
		} else {
			bf.readDir()
		}
	} else {
		bf.readFile()
	}
}

// walk遍历文件夹
func (bf *BasicFtp) walkDir() {
	c, err := bf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()
	walker := c.Walk(bf.Path)
	switch bf.Reg {
	case nil:
		bf.walkDirMode(walker)
	default:
		bf.walkDirModeReg(walker)
	}
}

func (bf *BasicFtp) walkDirMode(walker *ftp.Walker) {
	for walker.Next() {
		_, next := bf.walkDirModeBase(walker)
		if !next {
			continue
		}
		fmt.Println(walker.Path())
	}
}

func (bf *BasicFtp) walkDirModeReg(walker *ftp.Walker) {
	for walker.Next() {
		entry, next := bf.walkDirModeBase(walker)
		if !next {
			continue
		}
		if bf.Reg.FindString(entry.Name) != "" {
			fmt.Println(walker.Path())
		}
	}
}

func (bf *BasicFtp) walkDirModeBase(walker *ftp.Walker) (*ftp.Entry, bool) {
	entry := walker.Stat()
	switch bf.Mode {
	case "file":
		if entry.Type.String() == "folder" {
			return nil, false
		}
	case "dir":
		if entry.Type.String() != "folder" {
			return nil, false
		}
	}
	return entry, true
}

// 遍历文件目录
func (bf *BasicFtp) readDir() {
	c, err := bf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	entries, err := c.List(bf.Path)
	if err != nil {
		log.Fatal(err)
	}
	switch bf.Reg {
	case nil:
		bf.readDirMode(entries)
	default:
		bf.readDirModeReg(entries)
	}
}

func (bf *BasicFtp) readDirMode(entries []*ftp.Entry) {
	for _, entry := range entries {
		bf.readDirModeBase(entry)
	}
}

func (bf *BasicFtp) readDirModeReg(entries []*ftp.Entry) {
	for _, entry := range entries {
		if bf.Reg.FindString(entry.Name) == "" {
			continue
		}
		bf.readDirModeBase(entry)
	}
}

func (bf *BasicFtp) readDirModeBase(entry *ftp.Entry) {
	fType := ""
	switch bf.Mode {
	case "dir":
		if entry.Type.String() != "folder" {
			return
		}
		fType = "<DIR>"
	case "file":
		if entry.Type.String() == "folder" {
			return
		}
	default:
		if entry.Type.String() == "folder" {
			fType = "<DIR>"
		}
	}
	fmt.Printf("%s %5s %15d %s\n",
		entry.Time.Format("2006-01-02 15:04:05"),
		fType,
		entry.Size,
		entry.Name)
}

// 阅读单个文件
func (bf *BasicFtp) readFile() {
	c, err := bf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	res, err := c.Retr(bf.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()

	bs := bufio.NewScanner(res)
	switch bf.Reg {
	case nil:
		for bs.Scan() {
			fmt.Println(bs.Text())
		}
	default:
		for bs.Scan() {
			if bf.Reg.FindString(bs.Text()) != "" {
				fmt.Println(bs.Text())
			}
		}
	}
}
