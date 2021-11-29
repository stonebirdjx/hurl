// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: basic.go
// @Date: 2021/11/20 11:15
// @Desc: file protocol basic information
package file

import (
	"hurl/configs"
	"hurl/shares"
	"os"
	"regexp"
)

// 文件协议结构体信息
type BasicFile struct {
	Path  string         // 本地文件路径
	IsDir bool           // 是否是文件夹
	Walk  bool           // 是否walk遍历当前文件夹
	Mode  string         // walk类型必须取值all, dir, file之一
	Reg   *regexp.Regexp // 正则匹配规则
}

// 构建一个文件协议结构体指针
func NewBasicFiler(path string) (*BasicFile, error) {
	pathStat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	reg, err := shares.IfReg()
	if err != nil {
		return nil, err
	}

	return &BasicFile{
		Path:  path,
		IsDir: pathStat.IsDir(),
		Walk:  *configs.Walk,
		Mode:  *configs.Mode,
		Reg:   reg,
	}, nil
}

// 文件协议统一入口
func (b *BasicFile) Entrance() {
	if b.IsDir {
		b.folder() // 文件夹浏览
	} else {
		b.single() // 文件浏览
	}
}
