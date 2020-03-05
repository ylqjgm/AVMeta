package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// MakeRequest 发起请求
func MakeRequest(
	method, uri, proxy string,
	body io.Reader,
	header map[string]string,
	cookies []*http.Cookie) (
	data []byte,
	status int,
	err error) {
	// 构建请求客户端
	client := createHTTPClient(proxy)

	// 创建请求对象
	req, err := createRequest(method, uri, body, header, cookies)
	// 检查错误
	if err != nil {
		return nil, 0, err
	}

	// 执行请求
	res, err := client.Do(req)
	// 检查错误
	if err != nil {
		return nil, 0, err
	}

	// 获取请求状态码
	status = res.StatusCode
	// 读取请求内容
	data, err = ioutil.ReadAll(res.Body)
	// 关闭请求连接
	_ = res.Body.Close()

	return data, status, err
}

// GetResult 获取远程字节集数据
func GetResult(uri, proxy string, cookies []*http.Cookie) ([]byte, error) {
	// 头部定义
	header := make(map[string]string)
	// 加入头部信息
	header["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/68.0.3440.106 Safari/537.36"

	// 执行请求
	body, status, err := MakeRequest("GET", uri, proxy, nil, header, cookies)
	// 检查错误
	if err != nil {
		return nil, err
	}

	// 检查状态码
	if http.StatusBadRequest <= status {
		err = fmt.Errorf("%d", status)
	}

	return body, err
}

// GetRoot 获取远程节点数据
func GetRoot(uri, proxy string, cookies []*http.Cookie) (*goquery.Document, error) {
	// 获取远程字节集数据
	data, err := GetResult(uri, proxy, cookies)
	// 检查错误
	if err != nil {
		return nil, err
	}

	// 转换为节点数据
	root, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	// 检查错误
	if err != nil {
		return nil, err
	}

	return root, nil
}

// SavePhoto 远程图片下载
func SavePhoto(uri, savePath, proxy string, needConvert bool) error {
	// 创建路径
	err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	// 检查错误
	if err != nil {
		return err
	}

	// 读取远程字节集
	body, err := GetResult(uri, proxy, nil)
	// 检查错误
	if err != nil {
		return err
	}

	// 获取远程图片大小
	length := int64(len(body))
	// 检查大小
	if length == 0 || length < 1024 {
		return fmt.Errorf("远程图片不完整或小于1KB")
	}

	// 保存到本地
	err = saveFile(savePath, body, length)
	// 检查错误
	if err != nil {
		return err
	}

	// 是否需要转换
	if needConvert {
		// 转换为jpg
		err = ConvertJPG(savePath, fmt.Sprintf("%s.jpg", strings.TrimRight(path.Base(savePath), path.Ext(savePath))))
		// 检查
		if err != nil {
			return err
		}

		// 删除源文件
		return os.Remove(savePath)
	}

	return nil
}

// 创建http客户端
func createHTTPClient(proxy string) *http.Client {
	// 初始化
	transport := &http.Transport{
		/* #nosec */
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 如果有代理
	if proxy != "" {
		// 解析代理地址
		proxyURI := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxy)
		}
		// 加入代理
		transport.Proxy = proxyURI
	}

	// 返回客户端
	return &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}
}

// 创建请求对象
func createRequest(method, uri string, body io.Reader, header map[string]string, cookies []*http.Cookie) (*http.Request, error) {
	// 新建请求
	req, err := http.NewRequest(method, uri, body)
	// 检查错误
	if err != nil {
		return nil, err
	}

	// 循环头部信息
	for k, v := range header {
		// 设置头部
		req.Header.Set(k, v)
	}

	// 设置了cookie
	if len(cookies) > 0 {
		// 循环cookie
		for _, cookie := range cookies {
			// 加入cookie
			req.AddCookie(cookie)
		}
	}

	return req, err
}

// 保存字节集到本地
func saveFile(savePath string, data []byte, length int64) error {
	// 定义错误变量
	var err error

	// 创建路径
	_ = os.MkdirAll(path.Dir(savePath), os.ModePerm)

	// 创建空文件
	f, err := os.Create(savePath)
	// 检查错误
	if err != nil {
		return err
	}

	// 读取数据
	rc := bytes.NewReader(data)
	// 拷贝到指定路径
	_, err = io.Copy(f, rc)
	// 检查错误
	if err != nil {
		return err
	}

	// 关闭连接
	_ = f.Close()

	// 获取文件大小
	local := GetFileSize(savePath)

	// 检查文件一致性
	if length != local {
		// 删除已下载文件
		_ = os.Remove(savePath)
		// 设置错误信息
		err = fmt.Errorf("文件不完成, 下载失败")
	}

	return err
}
