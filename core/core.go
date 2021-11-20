// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: core.go
// @Date: 2021/11/20 10:49
// @Desc: core shunt layer
package core

import (
	"hurl/core/file"
	"log"
	"net/url"
)

// 文件协议消息处理者
// 传入类型 *url.URL
func FileHandle(u *url.URL) {
	path := u.Path
	basicFile, err := file.NewBasicFiler(path)
	if err != nil {
		log.Fatal(err)
	}
	basicFile.Entrance()
}
