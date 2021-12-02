// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: sftpupload.go
// @Date: 2021/11/28 9:06
// @Desc: sftp upload processor
package ftpsftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"hurl/configs"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// sftp上传路由
func (bsf *BasicSftp) Upload() {
	// 文件上传sftp必须是路径，path 以/结尾
	if !strings.HasSuffix(bsf.Path, "/") {
		log.Fatal("upload mode ftp path must end with /")
	}

	local := strings.TrimSpace(*configs.Upload)

	fileInfo, err := os.Stat(local)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		bsf.uploadDir(local) // 文件夹上传
	} else {
		bsf.uploadFile(fileInfo) // 单个文件上传
	}
}

// sftp 文件夹上传
func (bsf *BasicSftp) uploadDir(localDir string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer close(trChan)
		var err error
		switch bsf.Reg {
		case nil:
			err = filepath.Walk(localDir, bsf.visit)
		default:
			err = filepath.Walk(localDir, bsf.visitReg)
		}

		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	for i := 0; i < *configs.Currency; i++ {
		wg.Add(1)
		go func(i int) {
			bsf.uploadRangeFile(i, localDir)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// sftp多线程上传
func (bsf *BasicSftp) uploadRangeFile(i int, local string) {
	thread := "[sftp-upload-thread-" + strconv.Itoa(i) + "]:"
	c, err := bsf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	for tr := range trChan {
		start := float64(time.Now().UnixNano())
		relative := strings.TrimPrefix(filepath.ToSlash(tr.site), filepath.ToSlash(local))
		sftpFile := filepath.ToSlash(filepath.Join(bsf.Path, relative))
		sftpDir := sftpFile

		if tr.tp == configs.File {
			sftpDir = filepath.Dir(sftpDir)
		}
		sftpDir = filepath.ToSlash(sftpDir)

		// sftp 路径检查
		err = cmSftpPath(c, sftpDir)
		if err != nil {
			log.Fatal(err)
		}

		if tr.tp == configs.Dir {
			continue
		}

		err = uploadBase(c, sftpFile, tr.site)
		if err != nil {
			log.Fatal(err)
		}

		end := float64(time.Now().UnixNano())
		fmt.Printf("%s upload %s success totol-size:%d waste-time:%.2fms\n",
			thread,
			tr.site,
			tr.size,
			(end-start)/1e6)
	}
}

func (bsf *BasicSftp) visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	var tr transport
	tr.name = info.Name()
	tr.size = uint64(info.Size())
	if info.IsDir() {
		tr.tp = configs.Dir
	} else {
		tr.tp = configs.File
	}
	tr.site = path
	trChan <- tr
	return nil
}

func (bsf *BasicSftp) visitReg(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if bsf.Reg.FindString(info.Name()) == configs.EmptyString {
		return nil
	}

	err = bsf.visit(path, info, err)
	if err != nil {
		return err
	}
	return nil
}

// sftp 上传单个文件
func (bsf *BasicSftp) uploadFile(fileInfo os.FileInfo) {
	c, err := bsf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	start := float64(time.Now().UnixNano())

	// sftp 路径检查
	err = cmSftpPath(c, bsf.Path)
	if err != nil {
		log.Fatal(err)
	}

	// 上传处理
	local := strings.TrimSpace(*configs.Upload)
	ftpFile := filepath.ToSlash(filepath.Join(bsf.Path, fileInfo.Name()))
	err = uploadBase(c, ftpFile, local)
	if err != nil {
		log.Fatal(err)
	}

	end := float64(time.Now().UnixNano())
	fmt.Printf("sftp upload %s success totol-size:%d waste-time:%.2fms\n",
		local,
		fileInfo.Size(),
		(end-start)/1e6)
}

func uploadBase(c *sftp.Client, ftpFile, localFile string) error {
	lf, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer lf.Close()

	sf, err := c.OpenFile(ftpFile, os.O_CREATE|os.O_RDWR)
	if err != nil {
		return err
	}
	defer sf.Close()

	buff := make([]byte, *configs.ReadBytes)
	accept := 0
	for {
		n, err := lf.Read(buff)
		if n > 0 {
			_, err := sf.Write(buff[0:n])
			if err != nil {
				log.Fatal(err)
			}
			accept = accept + n
			fmt.Printf("uploading %s, send-byte:%d\r", localFile, accept)
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
	}
	return nil
}

// 检查sftp path 是否存在
// 不存在则创建
func cmSftpPath(c *sftp.Client, sftpPath string) error {
	mutex.Lock()
	defer mutex.Unlock()
	fileInfo, err := c.Stat(sftpPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := c.MkdirAll(sftpPath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !fileInfo.IsDir() {
		return fmt.Errorf("sftp path %s is not a dir", sftpPath)
	}
	return nil
}
