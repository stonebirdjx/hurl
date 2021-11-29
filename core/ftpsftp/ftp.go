// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: ftp.go
// @Date: 2021/11/23 22:00
// @Desc:
package ftpsftp

import (
	"github.com/jlaffaye/ftp"
	"time"
)

// ftp协议客户端登录
func (bf *BasicFtp) login() (*ftp.ServerConn, error) {
	c, err := ftp.Dial(bf.Host, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	err = c.Login(bf.User, bf.Password)
	if err != nil {
		return nil, err
	}

	return c, nil
}
