// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: file.go
// @Date: 2021/11/20 11:00
// @Desc: deal file protocol path is file
package file

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// 文件协议单个文件消息入口
// p: 文件路径
func (b *BasicFile) single() {
	f, err := os.Open(b.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bs := bufio.NewScanner(f)
	switch b.Reg {
	case nil:
		for bs.Scan() {
			fmt.Println(bs.Text())
		}
	default:
		for bs.Scan() {
			if b.Reg.FindString(bs.Text()) != "" {
				fmt.Println(bs.Text())
			}
		}
	}
}
