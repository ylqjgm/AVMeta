package capture

import (
	"bytes"
	/* #nosec */
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
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

const (
	// HEYZO 常量
	HEYZO = "HEYZO"
	// TOKYOHOT 常量
	TOKYOHOT = "東京熱"
	// FC2 常量
	FC2 = "FC2"

	// ZERO 数字0
	ZERO = 0
	// ONE 数字1
	ONE = 1
	// TWO 数字2
	TWO = 2
)

// ICapture 刮削器接口
type ICapture interface {
	// Fetch 获取数据
	Fetch(code string) error

	// GetURI 获取来源页面地址
	GetURI() string

	// GetNumber 获取刮削番号
	GetNumber() string

	// GetTitle 获取影片名称
	GetTitle() string
	// GetIntro 获取影片简介
	GetIntro() string
	// GetDirector 获取影片导演
	GetDirector() string
	// GetRelease 获取发行时间
	GetRelease() string
	// GetRuntime 获取影片时长
	GetRuntime() string
	// GetStudio 获取影片厂商
	GetStudio() string
	// GetSerise 获取影片系列
	GetSerise() string
	// GetTags 获取标签列表
	GetTags() []string
	// GetCover 获取封面图片
	GetCover() string
	// GetActors 获取演员列表
	GetActors() map[string]string
}

// MD5Verify md5验证
func MD5Verify(data []byte, source string) bool {
	/* #nosec */
	ret := md5.Sum(data)
	// 获取加密字符串
	md5Str := hex.EncodeToString(ret[:])

	// 返回比较结果
	return source == md5Str
}

// CheckDomainPrefix 检查域名最后的斜线
func CheckDomainPrefix(domain string) string {
	// 是否为空
	if domain == "" {
		return ""
	}

	// 获取最后一个字符
	last := domain[len(domain)-1:]
	// 如果是斜线
	if last == "/" {
		domain = domain[:len(domain)-1]
	}

	return domain
}

// IntroFilter 过滤简介
func IntroFilter(intro string) string {
	// 替换<br>
	intro = strings.ReplaceAll(intro, "<br>", "\n")
	intro = strings.ReplaceAll(intro, "<br/>", "\n")
	intro = strings.ReplaceAll(intro, "<br />", "\n")
	// 替换\r\n
	intro = strings.ReplaceAll(intro, "\r\n", "\n")
	// 替换\r
	intro = strings.ReplaceAll(intro, "\r", "\n")
	// 替换\n\n
	intro = strings.ReplaceAll(intro, "\n\n", "\n")

	// 清除多余空白
	return strings.TrimSpace(intro)
}

// ConvertJPG 转换图片为jpg格式
func ConvertJPG(sourceFile, newFile string) error {
	// 打开原文件
	f, err := os.Open(sourceFile)
	// 检查错误
	if err != nil {
		return err
	}
	// 关闭连接
	defer f.Close()

	// 图片解码
	src, _, err := image.Decode(f)
	// 检查错误
	if err != nil {
		return err
	}

	// 获取图片信息
	b := src.Bounds()

	// YCBCr
	img := src.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, 0, b.Max.X, b.Max.Y))

	// 新建并打开新图片
	cf, err := os.OpenFile(newFile, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	// 检查错误
	if err != nil {
		return err
	}
	// 关闭
	defer cf.Close()

	// 图片编码
	return jpeg.Encode(cf, img, &jpeg.Options{Quality: 100})
}

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
func SavePhoto(uri, savePath, proxy string) error {
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

	// 进行MD5验证
	if MD5Verify(body, "f591f3826a1085af5cdeeca250b2c97a") {
		return fmt.Errorf("远程图片为空图片或错误图片")
	}

	// 保存到本地
	return saveFile(savePath, body, length)
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

	// 获取文件信息
	info, err := os.Stat(savePath)
	// 检查错误
	if err != nil {
		return err
	}

	// 获取文件大小
	local := info.Size()

	// 检查文件一致性
	if length != local {
		// 删除已下载文件
		_ = os.Remove(savePath)
		// 设置错误信息
		err = fmt.Errorf("文件不完成, 下载失败")
	}

	return err
}
