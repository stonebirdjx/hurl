// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: file.go
// @Date: 2021/11/20 11:00
// @Desc: deal file protocol path is file
package file

import "fmt"

// 文件协议单个文件消息入口
// p: 文件路径
func (b *BasicFile) single() {
	fmt.Println("file")
}
