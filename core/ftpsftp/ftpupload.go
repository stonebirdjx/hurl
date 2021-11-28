// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: ftpupload.go
// @Date: 2021/11/28 18:07
// @Desc: ftp upload processor
package ftpsftp

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"hurl/configs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func (bf *BasicFtp) Upload() {
	if !strings.HasSuffix(bf.Path, "/") {
		log.Fatal("upload mode ftp path must end with /")
	}
	local := strings.TrimSpace(*configs.Upload)
	fileInfo, err := os.Stat(local)
	if err != nil {
		log.Fatal(err)
	}
	if fileInfo.IsDir() {
		bf.uploadDir(local) // deal with dir
	} else {
		bf.uploadFile(fileInfo) // deal with file
	}
}

func (bf *BasicFtp) uploadDir(local string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := filepath.Walk(local, bf.visit)
		if err != nil {
			log.Fatal(err)
		}
		close(trChan)
		wg.Done()
	}()
	for i := 0; i < *configs.Currency; i++ {
		wg.Add(1)
		go func(i int) {
			bf.uploadRangeFile(i, local)
			wg.Done()
		}(i)
	}
	wg.Wait()

}

func (bf *BasicFtp) uploadRangeFile(i int, local string) {
	thread := "[ftp-upload-thread-" + fmt.Sprint(i) + "]:"

	c, err := bf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	for tr := range trChan {
		start := float64(time.Now().UnixNano())
		relative := strings.TrimPrefix(tr.site, local)
		ftpFile := filepath.Join(bf.Path, relative)
		ftpDir := filepath.ToSlash(ftpFile)
		if tr.tp == "file" {
			ftpDir = filepath.Dir(ftpDir)
		}
		ftpDir = filepath.ToSlash(ftpDir)
		err := cmFtpPath(c, ftpDir)
		if err != nil {
			log.Fatal(err)
		}
		if tr.tp == "dir" {
			continue
		}
		err = ftpUploadBase(c, ftpFile, tr.site)
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

func (bf *BasicFtp) visit(fp string, info os.FileInfo, err error) error {
	if err != nil {
		log.Fatalln(err)
	}
	fp = filepath.ToSlash(fp)
	if bf.Reg != nil {
		if bf.Reg.FindString(info.Name()) == "" {
			return nil
		}
	}
	var tr transport
	tr.name = info.Name()
	tr.size = uint64(info.Size())
	if info.IsDir() {
		tr.tp = "dir"
	} else {
		tr.tp = "file"
	}
	tr.site = fp
	trChan <- tr
	return nil
}

func (bf *BasicFtp) uploadFile(fileInfo os.FileInfo) {
	start := float64(time.Now().UnixNano())
	c, err := bf.login()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	err = cmFtpPath(c, bf.Path)
	if err != nil {
		log.Fatal(err)
	}

	localFile := strings.TrimSpace(*configs.Upload)
	ftpFile := filepath.ToSlash(filepath.Join(bf.Path, fileInfo.Name()))
	err = ftpUploadBase(c, ftpFile, localFile)
	if err != nil {
		log.Fatal(err)
	}
	end := float64(time.Now().UnixNano())
	fmt.Printf("ftp upload %s success totol-size:%d waste-time:%.2fms\n",
		localFile,
		fileInfo.Size(),
		(end-start)/1e6)
}

func ftpUploadBase(c *ftp.ServerConn, ftpFile, localFile string) error {
	file, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer file.Close()

	err = c.Stor(ftpFile, file)
	if err != nil {
		return err
	}
	return nil
}

func cmFtpPath(c *ftp.ServerConn, path string) error {
	mutex.Lock()
	defer mutex.Unlock()
	currentDir, err := c.CurrentDir()
	if err != nil {
		return err
	}

	paths := strings.Split(path, "/")
	if filepath.IsAbs(path) {
		err = c.ChangeDir("/")
		if err != nil {
			return err
		}
	}
	for _, p := range paths {
		err := checkPath(c, p)
		if err != nil {
			return err
		}
	}
	err = c.ChangeDir(currentDir)
	if err != nil {
		return err
	}
	return nil
}

func checkPath(c *ftp.ServerConn, path string) error {
	if path == "" {
		return nil
	}
	err := c.ChangeDir(path)
	if err != nil {
		mkdirError := c.MakeDir(path)
		if mkdirError != nil {
			return err
		}
		changeError := c.ChangeDir(path)
		if changeError != nil {
			return err
		}
	}
	return nil
}
