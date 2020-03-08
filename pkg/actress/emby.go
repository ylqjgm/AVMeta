package actress

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"
)

// Emby Emby媒体库结构体
type Emby struct {
	// Emby访问地址
	hostURL string
	// API
	apiKey string
}

// 用户结构
type embyPerson struct {
	Name      string      `json:"Name"`
	ID        string      `json:"Id"`
	ImageTags embyPrimary `json:"ImageTags"`
}

// 封面结构
type embyPrimary struct {
	Primary string `json:"Primary"`
}

// NewEmby 返回一个被初始化的 Emby 对象
//
// hostURL 字符串参数，传入 Emby 媒体库访问地址，
// apiKey 字符串参数，传入 Emby 的 API Key。
func NewEmby(hostURL, apiKey string) *Emby {
	return &Emby{
		hostURL: hostURL,
		apiKey:  apiKey,
	}
}

// Actor 单个女优头像入库
//
// name 字符串参数，传入女优姓名，必须与 Emby 中一致才会入库，
// face 字符串参数，传入女优头像图片本地路径。
func (emby *Emby) Actor(name, face string) error {
	// 检查传入数据
	if name == "" || face == "" {
		return fmt.Errorf("演员名字或头像路径不能为空")
	}

	// 获取演员信息
	per, err := emby.getPerson(name)
	// 检查错误
	if err != nil {
		return err
	}

	// 检查是否已经有了
	if per.ImageTags.Primary != "" {
		return nil
	}

	return emby.uploadImage(per.ID, face)
}

// 本地上传演员头像
func (emby *Emby) uploadImage(id, face string) error {
	// 图片编码
	body, err := emby.base64Encode(face)
	// 检查
	if err != nil {
		return err
	}

	// 组合地址
	uri := fmt.Sprintf("emby/Items/%s/Images/Primary", id)
	// 提交请求
	_, err = emby.makeRequest("POST", uri, body)

	return err
}

// 获取演员信息
func (emby *Emby) getPerson(name string) (*embyPerson, error) {
	// 发起请求
	raw, err := emby.makeRequest("GET", fmt.Sprintf("Persons/%s", url.PathEscape(name)), "")
	// 检查错误
	if err != nil {
		return nil, err
	}
	// 用户对象
	var per embyPerson
	// 将json解析到结构体中
	err = json.Unmarshal(raw, &per)

	return &per, err
}

// 发起请求
func (emby *Emby) makeRequest(method, uri, body string) ([]byte, error) {
	// 组合路径
	uri = fmt.Sprintf("%s/%s", emby.hostURL, uri)
	// 头部map
	header := make(map[string]string)
	// 设置API key
	header["X-Emby-Token"] = emby.apiKey

	// 发起请求
	data, status, err := util.MakeRequest(method, uri, "", strings.NewReader(body), header, nil)

	// 检查状态码
	if http.StatusOK != status && http.StatusNoContent != status && err != nil {
		err = fmt.Errorf("%d", status)
	}

	return data, err
}

// 文件Base64编码
func (emby *Emby) base64Encode(file string) (string, error) {
	// 检查错误
	f, err := os.Open(file)
	// 如果出错
	if err != nil {
		return "", err
	}
	// 关闭
	defer f.Close()

	// 初始化byte
	buff := make([]byte, 500000)
	// 读取文件
	n, err := f.Read(buff)
	// 检查错误
	if err != nil {
		return "", err
	}

	// Base64编码
	source := base64.StdEncoding.EncodeToString(buff[:n])

	return source, nil
}
