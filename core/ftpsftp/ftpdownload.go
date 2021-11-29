// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: ftpdownload.go
// @Date: 2021/11/28 19:24
// @Desc: ftp download processor
package ftpsftp

import (
	"bufio"
	"fmt"
	"github.com/jlaffaye/ftp"
	"hurl/configs"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ftp协议下载时入口
func (bf *BasicFtp) Download() {
	local := strings.TrimSpace(*configs.Download)
	// 本地路径判断
	downloadPathIsDir(local)

	if strings.HasSuffix(bf.Path, "/") {
		bf.downloadDir(local) //下载的是目录
	} else {
		bf.downloadFile(local) //下载的是文件
	}
}

// ftp 下载目录处理
func (bf *BasicFtp) downloadDir(local string) {
	c, err := bf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		walker := c.Walk(bf.Path)
		for walker.Next() {
			if tr, ok := bf.toChan(walker); ok {
				trChan <- tr
			}
		}
		close(trChan)
		wg.Done()
	}()

	for i := 0; i < *configs.Currency; i++ {
		wg.Add(1)
		go func(i int) {
			bf.downloadRangFile(i, local)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// chan消息处理
func (bf *BasicFtp) toChan(walker *ftp.Walker) (transport, bool) {
	entry := walker.Stat()
	if bf.Reg != nil && bf.Reg.FindString(entry.Name) == configs.EmptyString {
		return transport{}, false
	}

	var tr transport
	tr.name = entry.Name
	tr.tp = entry.Type.String()
	tr.size = entry.Size
	tr.site = walker.Path()

	filePath := tr.site
	if strings.HasPrefix(filePath, "/") {
		filePath = strings.TrimLeft(filePath, "/")
	}

	tmp := ""
	if strings.HasPrefix(bf.Path, "/") {
		tmp = strings.TrimLeft(bf.Path, "/")
	} else if strings.HasPrefix(bf.Path, "./") {
		tmp = strings.TrimLeft(bf.Path, "./")
	}

	tr.relative = strings.TrimPrefix(filePath, tmp)
	return tr, true
}

// 多携程下载ftp文件
func (bf *BasicFtp) downloadRangFile(i int, local string) {
	thread := "[ftp-download-thread-" + fmt.Sprint(i) + "]:"
	c, err := bf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	for tr := range trChan {
		start := float64(time.Now().UnixNano())
		localPath := filepath.Join(local, tr.relative)
		dir := localPath
		if tr.tp != configs.Dir {
			dir = filepath.Dir(localPath)
		}

		// 检查本地目录
		err := cmLocalDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		// 目录创建不下载
		if tr.tp == configs.Dir {
			continue
		}

		// 下载文件
		err = ftpDownloadBase(c, tr.site, localPath)
		if err != nil {
			log.Fatal(err)
		}

		end := float64(time.Now().UnixNano())
		fmt.Printf("%s download %s success totol-size:%d waste-time:%.2fms\n",
			thread,
			tr.site,
			tr.size,
			(end-start)/1e6)
	}
}

// 单个文件下载入口
func (bf *BasicFtp) downloadFile(local string) {
	c, err := bf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	start := float64(time.Now().UnixNano())
	entries, err := c.List(bf.Path)
	if err != nil {
		log.Fatal(err)
	}

	// clean code 编码防止魔鬼数字
	l := 1
	if len(entries) != l {
		log.Fatal("ftp path dir must end with /")
	}

	localPath := filepath.Join(local, entries[0].Name)

	// 下载文件
	err = ftpDownloadBase(c, localPath, bf.Path)
	if err != nil {
		log.Fatal(err)
	}

	end := float64(time.Now().UnixNano())
	fmt.Printf("download %s success totol-size:%d waste-time:%.2fms\n",
		local,
		entries[0].Size,
		(end-start)/1e6)
}

// 下载基本方法处理
// c ftp客户端
// ftpFile:ftp上的文件
// localFile:本地存放的文件
func ftpDownloadBase(c *ftp.ServerConn, ftpFile, localFile string) error {
	resp, err := c.Retr(ftpFile)
	if err != nil {
		return err
	}
	defer resp.Close()

	f, err := os.OpenFile(localFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	br := bufio.NewReader(resp)
	buff := make([]byte, *configs.ReadBytes)
	accept := 0
	for {
		n, err := br.Read(buff)
		if n > 0 {
			_, err = f.Write(buff[0:n])
			if err != nil {
				return err
			}
			accept = accept + n
			fmt.Printf("downloading %s, accept-byte:%d\r", ftpFile, accept)
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
