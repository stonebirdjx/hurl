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

func (bf *BasicFtp) Download() {
	local := strings.TrimSpace(*configs.Download)
	localInfo, err := os.Stat(local)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(local, 0644)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	} else if !localInfo.IsDir() {
		log.Fatal("local path is not a dir")
	}

	if strings.HasSuffix(bf.Path, "/") {
		bf.downloadDir(local)
	} else {
		bf.downloadFile(local)
	}
}

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

func (bf *BasicFtp) toChan(walker *ftp.Walker) (transport, bool) {
	entry := walker.Stat()
	if bf.Reg != nil {
		if bf.Reg.FindString(entry.Name) == "" {
			return transport{}, false
		}
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
		if tr.tp != "dir" {
			dir = filepath.Dir(localPath)
		}
		err := cmLocalDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		if tr.tp == "dir" {
			continue
		}

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

func (bf *BasicFtp) downloadFile(local string) {
	start := float64(time.Now().UnixNano())
	c, err := bf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	entries, err := c.List(bf.Path)
	if err != nil {
		log.Fatal(err)
	}

	if len(entries) != 1 {
		log.Fatal("ftp path dir must end with /")
	}

	localPath := filepath.Join(local, entries[0].Name)
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
