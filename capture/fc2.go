package capture

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FC2Capture fc2刮削器
type FC2Capture struct {
	Proxy       string            // 代理设置
	uri         string            // 页面地址
	code        string            // 临时番号
	number      string            // 最终番号
	fc2Root     *goquery.Document // fc2根节点
	fc2clubRoot *goquery.Document // fc2club根节点
}

// fc2标签json结构
type fc2tags struct {
	Tags []fc2tag `json:"tags"`
}

// fc2标签内容结构
type fc2tag struct {
	Tag string `json:"tag"`
}

// Fetch 刮削
func (s *FC2Capture) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 过滤番号
	r, _ := regexp.Compile(`[0-9]{6,7}`)
	// 获取临时番号
	s.code = r.FindString(code)
	// 组合fc2地址
	fc2uri := fmt.Sprintf("https://adult.contents.fc2.com/article/%s/", s.code)
	// 组合fc2club地址
	fc2cluburi := fmt.Sprintf("https://fc2club.com/html/FC2-%s.html", s.code)

	// 打开fc2
	fc2Root, err := GetRoot(fc2uri, s.Proxy, nil)
	// 检查错误
	if nil != err {
		return err
	}

	// 打开fc2club
	fc2clubRoot, err := GetRoot(fc2cluburi, s.Proxy, nil)
	// 检查错误
	if nil != err {
		return err
	}

	// 设置页面地址
	s.uri = fc2uri
	// 设置fc2根节点
	s.fc2Root = fc2Root
	// 设置fc2club根节点
	s.fc2clubRoot = fc2clubRoot

	return nil
}

// GetTitle 获取名称
func (s *FC2Capture) GetTitle() string {
	// 获取标题
	title := s.fc2Root.Find(`.items_article_headerInfo h3`).Text()
	// 检查
	if "" == title {
		title = s.fc2clubRoot.Find(`.main h3`).Text()
	}

	return title
}

// GetIntro 获取简介
func (s *FC2Capture) GetIntro() string {
	return ""
}

// GetDirector 获取导演
func (s *FC2Capture) GetDirector() string {
	// 获取导演
	director := s.fc2Root.Find(`.items_article_headerInfo li:nth-child(3) a`).Text()
	// 检查
	if "" == director {
		director = s.fc2clubRoot.Find(`.main h5:nth-child(5) a:nth-child(2)`).Text()
	}

	return director
}

// GetRelease 发行时间
func (s *FC2Capture) GetRelease() string {
	return strings.ReplaceAll(strings.ReplaceAll(s.fc2Root.Find(`.items_article_Releasedate p`).Text(), "上架时间 :", ""), "販売日 :", "")
}

// GetRuntime 获取时长
func (s *FC2Capture) GetRuntime() string {
	return "0"
}

// GetStudio 获取厂商
func (s *FC2Capture) GetStudio() string {
	return "FC2"
}

// GetSerise 获取系列
func (s *FC2Capture) GetSerise() string {
	return "FC2"
}

// GetTags 获取标签
func (s *FC2Capture) GetTags() []string {
	// 组合地址
	uri := fmt.Sprintf("http://adult.contents.fc2.com/api/v4/article/%s/tag?", s.code)

	// 读取远程数据
	data, err := GetResult(uri, s.Proxy, nil)
	// 检查
	if nil != err {
		return nil
	}

	// 读取内容
	body, err := ioutil.ReadAll(bytes.NewReader(data))
	// 检查错误
	if nil != err {
		return nil
	}

	// json
	var tagsJSON fc2tags

	// 解析json
	err = json.Unmarshal(body, &tagsJSON)
	// 检查
	if nil != err {
		return nil
	}

	// 定义数组
	var tags []string

	// 循环标签
	for _, tag := range tagsJSON.Tags {
		tags = append(tags, strings.TrimSpace(tag.Tag))
	}

	return tags
}

// GetFanart 获取图片
func (s *FC2Capture) GetFanart() string {
	// 获取图片
	fanart, _ := s.fc2clubRoot.Find(`.slides li:nth-child(1) img`).Attr("src")
	// 检查
	if "" == fanart {
		return ""
	}
	// 组合地址
	return fmt.Sprintf("https://fc2club.com%s", fanart)
}

// GetActors 获取演员
func (s *FC2Capture) GetActors() map[string]string {
	return nil
}

// GetURI 获取页面地址
func (s *FC2Capture) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *FC2Capture) GetNumber() string {
	return s.number
}
