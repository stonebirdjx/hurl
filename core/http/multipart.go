// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: multipart.go
// @Date: 2021/11/30 21:41
// @Desc: http multipart deal
package http

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"hurl/configs"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// http multipart
func (bh *BasicHttp) multipart() {
	data := strings.TrimSpace(*configs.Data)
	if !gjson.Valid(data) {
		log.Fatal("-multipart mode -data must json text")
	}

	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		if data != configs.EmptyString {
			var mp map[string]interface{}
			err := json.Unmarshal([]byte(data), &mp)
			if err != nil {
				log.Fatal(err)
			}
			writeField(mp, m)
		}

		filesArgs := strings.TrimSpace(*configs.HFile)
		if filesArgs != configs.EmptyString {
			files := strings.Split(filesArgs, ",")
			for index, file := range files {
				file = strings.TrimSpace(file)
				createFormFile(index, file, m)
			}
		}
	}()

	req, err := http.NewRequest(bh.Method, bh.Url, r)
	if err != nil {
		log.Fatal(err)
	}

	setHeader(req)
	req.Header.Set(contentType, m.FormDataContentType())
	bh.do(req)
}

// multipart write field
func writeField(mp map[string]interface{}, m *multipart.Writer) {
	for k, v := range mp {
		err := m.WriteField(k, fmt.Sprint(v))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createFormFile(index int, file string, m *multipart.Writer) {
	fi, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	part, err := m.CreateFormFile("file"+strconv.Itoa(index), fi.Name())
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	buff := make([]byte, *configs.ReadBytes)
	for {
		n, err := f.Read(buff)
		if n > 0 {
			_, err := part.Write(buff[0:n])
			if err != nil {
				log.Fatal(err)
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
	}
}
