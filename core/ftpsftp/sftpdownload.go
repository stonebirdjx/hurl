// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: sftpdownload.go
// @Date: 2021/11/28 14:45
// @Desc: sftp download processor
package ftpsftp

import (
	"fmt"
	"github.com/kr/fs"
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

// sftp下载时入口
func (bsf *BasicSftp) Download() {
	local := strings.TrimSpace(*configs.Download)
	downloadPathIsDir(local)

	if strings.HasSuffix(bsf.Path, "/") {
		bsf.downloadDir(local) // sftp 下载文件夹
	} else {
		bsf.downloadFile(local) // sftp 下载单个文件
	}
}

// 下载文件夹
func (bsf *BasicSftp) downloadDir(local string) {
	c, err := bsf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	spInfo, err := c.Stat(bsf.Path)
	if err != nil {
		log.Fatal(err)
	}

	if !spInfo.IsDir() {
		log.Fatal("ftp path is not a dir")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		walker := c.Walk(bsf.Path)
		for walker.Step() {
			if tr, ok := bsf.toChan(walker); ok {
				trChan <- tr
			}
		}
		close(trChan)
		wg.Done()
	}()

	for i := 0; i < *configs.Currency; i++ {
		wg.Add(1)
		go func(i int) {
			bsf.downloadRangFile(i, local)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// 往channel通道传输消息
func (bsf *BasicSftp) toChan(walker *fs.Walker) (transport, bool) {
	fileInfo := walker.Stat()
	if bsf.Reg != nil && bsf.Reg.FindString(fileInfo.Name()) == configs.EmptyString {
		return transport{}, false
	}

	var tr transport
	tr.name = fileInfo.Name()
	if fileInfo.IsDir() {
		tr.tp = configs.Dir
	} else {
		tr.tp = configs.File
	}

	tr.size = uint64(fileInfo.Size())
	tr.site = walker.Path()
	filePath := tr.site

	if strings.HasPrefix(filePath, "/") {
		filePath = strings.TrimLeft(filePath, "/")
	}

	tmp := ""
	if strings.HasPrefix(bsf.Path, "/") {
		tmp = strings.TrimLeft(bsf.Path, "/")
	} else if strings.HasPrefix(bsf.Path, "./") {
		tmp = strings.TrimLeft(bsf.Path, "./")
	}

	tr.relative = strings.TrimPrefix(filePath, tmp)
	return tr, true
}

func (bsf *BasicSftp) downloadRangFile(i int, local string) {
	thread := "[sftp-download-thread-" + strconv.Itoa(i) + "]:"
	c, err := bsf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	for tr := range trChan {
		start := float64(time.Now().UnixNano())
		localPath := filepath.Join(local, tr.relative)
		dir := localPath
		if tr.tp != configs.Dir {
			dir = filepath.Dir(localPath)
		}

		err = cmLocalDir(dir)
		if err != nil {
			log.Fatalf("%s check dir err %s", thread, err)
		}

		if tr.tp == configs.Dir {
			continue
		}

		err = downloadBase(c, tr.site, localPath)
		if err != nil {
			log.Fatalf("%s %s\n", thread, err)
		}

		end := float64(time.Now().UnixNano())
		fmt.Printf("%s download %s success totol-size:%d waste-time:%.2fms\n",
			thread,
			tr.site,
			tr.size,
			(end-start)/1e6)
	}
}

// 下载单个文件
func (bsf *BasicSftp) downloadFile(local string) {
	start := float64(time.Now().UnixNano())
	c, err := bsf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	fileInfo, err := c.Stat(bsf.Path)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		log.Fatal("sftp path err or is not file")
	}

	localFile := filepath.ToSlash(filepath.Join(local, fileInfo.Name()))

	err = downloadBase(c, bsf.Path, localFile)
	if err != nil {
		log.Fatal(err)
	}

	end := float64(time.Now().UnixNano())
	fmt.Printf("download %s success totol-size:%d waste-time:%.2fms\n",
		local,
		fileInfo.Size(),
		(end-start)/1e6)
}

// 下载基础方法
func downloadBase(c *sftp.Client, sftpFile, localFile string) error {
	sf, err := c.Open(sftpFile)
	if err != nil {
		return err
	}
	defer sf.Close()

	lf, err := os.OpenFile(localFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil
	}
	defer lf.Close()

	buff := make([]byte, *configs.ReadBytes)
	accept := 0
	for {
		n, err := sf.Read(buff)
		if n > 0 {
			_, err := lf.Write(buff[0:n])
			if err != nil {
				return nil
			}
			accept = accept + n
			fmt.Printf("downloading %s, accept-byte:%d\r", sftpFile, accept)
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
