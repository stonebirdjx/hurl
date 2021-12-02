// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: request.go
// @Date: 2021/12/2 11:03
// @Desc: other request
package http

import (
	"hurl/configs"
	"log"
	"net/http"
	"strings"
)

func (bh *BasicHttp) request() {
	data := strings.TrimSpace(*configs.Data)
	req, err := http.NewRequest(bh.Method, bh.Url, strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	setHeader(req)
	bh.do(req)
}
