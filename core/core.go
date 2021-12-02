// Copyright (c) 2021 hu. All rights reserved.
// @Author: stonebirdjx
// @Email: 1245863260@qq.com, g1245863260@gmail.com
// @File: core.go
// @Date: 2021/11/20 10:49
// @Desc: core shunt layer
package core

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/lucas-clemente/quic-go/http3"
	"golang.org/x/net/http2"
	"hurl/configs"
	"hurl/core/file"
	"hurl/core/ftpsftp"
	stbHttp "hurl/core/http"
	"hurl/shares"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// 文件协议消息处理者
// 传入类型 *url.URL
func FileHandle(u *url.URL) {
	path := u.Path
	basicFile, err := file.NewBasicFiler(path)
	if err != nil {
		log.Fatal(err)
	}
	basicFile.Entrance()
}

// sftp和ftp消息处理者
// 传入类型 *url.URL
func FtpSftpHandle(u *url.URL) {
	userName := u.User.Username()
	if userName == configs.EmptyString {
		userName = *configs.User
		if userName == configs.EmptyString {
			log.Fatalf("can not get ftp user")
		}
	}

	passWord, ok := u.User.Password()
	if !ok {
		passWord = *configs.Password
	}

	path := u.Path
	// 以//开头表绝对路径，否则是相对路径
	if !strings.HasPrefix(path, "//") {
		path = "." + path
	}

	// 判断是否使用正则表达式
	reg, err := shares.IfReg()
	if err != nil {
		log.Fatal(err)
	}

	var api ftpsftp.BasicApi
	basicStruct := ftpsftp.BasicStruct{
		Path:     path,
		Host:     u.Host,
		User:     userName,
		Password: passWord,
		Walk:     *configs.Walk,
		Mode:     *configs.Mode,
		Reg:      reg,
	}

	switch u.Scheme {
	case configs.FtpScheme:
		api = &ftpsftp.BasicFtp{
			BasicStruct: basicStruct,
		}
	case configs.SftpScheme:
		api = &ftpsftp.BasicSftp{
			BasicStruct: basicStruct,
		}
	default:
		log.Fatal("scheme is not ftp or sftp")
	}

	switch {
	case strings.TrimSpace(*configs.Download) != configs.EmptyString:
		api.Download()
	case strings.TrimSpace(*configs.Upload) != configs.EmptyString:
		api.Upload()
	default:
		api.Read()
	}
}

// http 协议消息处理者
func HttpHandle(u *url.URL) {
	client := &http.Client{}
	if u.Scheme == configs.HttpsScheme {
		httpsCrtCheck(client)
	}

	// 是否打开重定向
	if !*configs.HL {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // or maybe the error from the request
		}
	}

	// 请求方式
	method := strings.TrimSpace(*configs.Method)
	if *configs.HI {
		method = "HEAD"
	}

	bh := stbHttp.BasicHttp{
		Client: client,
		Method: method,
		Url:    u.String(),
	}

	bh.Entrance()
}

// https协议证书判断
func httpsCrtCheck(client *http.Client) {
	crtFile := strings.TrimSpace(*configs.Crt)
	switch crtFile {
	case configs.EmptyString:
		//跳过证书验证
		switch {
		case *configs.Http3:
			client.Transport = &http3.RoundTripper{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		case *configs.Http2:
			client.Transport = &http2.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		default:
			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		}
	default:
		// 自带证书验证
		caCert, err := ioutil.ReadFile(crtFile)
		if err != nil {
			log.Fatalf("Reading server certificate: %s", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Create TLS configuration with the certificate of the server
		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		switch {
		case *configs.Http3:
			client.Transport = &http3.RoundTripper{
				TLSClientConfig: tlsConfig,
			}
		case *configs.Http2:
			client.Transport = &http2.Transport{
				TLSClientConfig: tlsConfig,
			}
		default:
			client.Transport = &http.Transport{
				TLSClientConfig: tlsConfig,
			}
		}
	}
}
