// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: share.go
// @Date: 2021/11/21 10:52
// @Desc: some share function
package shares

import (
	"hurl/configs"
	"regexp"
)

// 判断是否启用正则公共方法
func IfReg() (*regexp.Regexp, error) {
	var re *regexp.Regexp
	var err error
	if *configs.Regexp != "" {
		re, err = regexp.Compile(*configs.Regexp)
	}
	return re, err
}
