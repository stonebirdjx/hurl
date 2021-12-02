// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: basic.go
// @Date: 2021/11/30 21:25
// @Desc: http protocol basic information
package http

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"hurl/configs"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	contentType    = "Content-Type"
	headMethod     = "HEAD"
	userAgent      = "User-Agent"
	userAgentValue = "hurl/@jx"
)

type BasicHttp struct {
	Client *http.Client
	Method string
	Url    string
}

func (bh *BasicHttp) Entrance() {
	if *configs.MultiPart {
		bh.multipart()
	} else {
		bh.request()
	}
}

func (bh *BasicHttp) do(req *http.Request) {
	res, err := bh.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	switch bh.Method {
	case headMethod:
		headRequest(res)
	default:
		otherRequest(res)
	}
}

// head 请求
func headRequest(res *http.Response) {
	fmt.Printf("%s %s\n", res.Proto, res.Status)
	for k, v := range res.Header {
		fmt.Printf("%s: %s\n", k, v)
	}
}

// 其他请求
func otherRequest(res *http.Response) {
	// i 是否开启
	if *configs.Hi {
		headRequest(res)
		fmt.Println() //留一行空行
	}

	saveName := strings.TrimSpace(*configs.Ho)
	if saveName != configs.EmptyString {
		saveResContent(saveName, res)
	} else {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}
}

// 保存网络文件
func saveResContent(saveName string, res *http.Response) {
	f, err := os.OpenFile(saveName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	contentLength := res.ContentLength
	start := time.Now().UnixNano()
	buffer := make([]byte, *configs.ReadBytes)
	accept := 0
	for {
		n, err := res.Body.Read(buffer)
		if n > 0 {
			_, err := f.Write(buffer[0:n])
			if err != nil {
				log.Fatal(err)
			}
			accept = accept + n
		}

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err, "accept-bytes:", accept)
			}
		}
		fmt.Printf("download accpet-byte:%d\r",
			accept,
		)
	}
	end := time.Now().UnixNano()
	fmt.Printf("totol:%d download:%d percentage:%.2f%% waste-time:%.2fms\r",
		contentLength,
		accept,
		float64(accept)/float64(contentLength)*100,
		float64(end-start)/1e6,
	)
}

// header 设置
func setHeader(req *http.Request) {
	headers := strings.TrimSpace(*configs.Headers)
	if headers != configs.EmptyString {
		if !gjson.Valid(headers) {
			log.Fatal("-headers must json text")
		}

		var mp map[string]interface{}
		err := json.Unmarshal([]byte(headers), &mp)
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := mp[userAgent]; !ok {
			req.Header.Set(userAgent, userAgentValue)
		}

		for k, v := range mp {
			req.Header.Set(k, fmt.Sprint(v))
		}

	} else {
		req.Header.Set(userAgent, userAgentValue)
	}
}
