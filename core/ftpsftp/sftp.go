// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: sftp.go
// @Date: 2021/11/23 22:00
// @Desc: sftp message processor
package ftpsftp

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"time"
)

// sftp 客户端
func (bsf *BasicSftp) login() (*sftp.Client, error) {
	config := ssh.ClientConfig{
		User:            bsf.User,
		Auth:            []ssh.AuthMethod{ssh.Password(bsf.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	sshClient, err := ssh.Dial("tcp", bsf.Host, &config)
	if err != nil {
		return nil, err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}
	return sftpClient, nil
}
