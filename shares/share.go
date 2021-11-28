// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: share.go
// @Date: 2021/11/21 10:52
// @Desc: some share function
package shares

import (
	"fmt"
	"hurl/configs"
	"io/fs"
	"regexp"
)

func PrintFileInfo(fileInfos []fs.FileInfo) {
	for _, fi := range fileInfos {
		printBase(fi)
	}
}

func PrintFileInfoReg(fileInfos []fs.FileInfo, reg *regexp.Regexp) {
	for _, fi := range fileInfos {
		if reg.FindString(fi.Name()) == "" {
			continue
		}
		printBase(fi)
	}
}

func printBase(fi fs.FileInfo) {
	fType := ""
	if fi.IsDir() {
		fType = "<DIR>"
	}
	fmt.Printf("%s %5s %15d %s\n",
		fi.ModTime().Format("2006-01-02 15:04:05"),
		fType,
		fi.Size(),
		fi.Name(),
	)
}

// 判断是否启用正则公共方法
func IfReg() (*regexp.Regexp, error) {
	var re *regexp.Regexp
	var err error
	if *configs.Regexp != "" {
		re, err = regexp.Compile(*configs.Regexp)
	}
	return re, err
}
