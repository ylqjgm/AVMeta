package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

// FC2Scraper fc2网站刮削器
type FC2Scraper struct {
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

// NewFC2Scraper 返回一个被初始化的fc2刮削对象
//
// proxy 字符串参数，传入代理信息
func NewFC2Scraper(proxy string) *FC2Scraper {
	return &FC2Scraper{Proxy: proxy}
}

// Fetch 刮削
func (s *FC2Scraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 过滤番号
	r := regexp.MustCompile(`[0-9]{6,7}`)
	// 获取临时番号
	s.code = r.FindString(code)
	// 组合fc2地址
	fc2uri := fmt.Sprintf("https://adult.contents.fc2.com/article/%s/", s.code)
	// 组合fc2club地址
	fc2cluburi := fmt.Sprintf("https://fc2club.net/html/FC2-%s.html", s.code)

	// 打开fc2
	fc2Root, err := util.GetRoot(fc2uri, s.Proxy, nil)
	// 检查错误
	if err != nil {
		return err
	}

	// 打开fc2club
	fc2clubRoot, err := util.GetRoot(fc2cluburi, s.Proxy, nil)
	// 检查错误
	if err != nil {
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
func (s *FC2Scraper) GetTitle() string {
	// 获取标题
	title := s.fc2Root.Find(`.items_article_headerInfo h3`).Text()
	// 检查
	if title == "" {
		title = s.fc2clubRoot.Find(`.main h3`).Text()
	}

	return title
}

// GetIntro 获取简介
func (s *FC2Scraper) GetIntro() string {
	return ""
}

// GetDirector 获取导演
func (s *FC2Scraper) GetDirector() string {
	// 获取导演
	director := s.fc2Root.Find(`.items_article_headerInfo li:nth-child(3) a`).Text()
	// 检查
	if director == "" {
		director = s.fc2clubRoot.Find(`.main h5:nth-child(5) a:nth-child(2)`).Text()
	}

	return director
}

// GetRelease 发行时间
func (s *FC2Scraper) GetRelease() string {
	return strings.ReplaceAll(strings.ReplaceAll(s.fc2Root.Find(`.items_article_Releasedate p`).Text(), "上架时间 :", ""), "販売日 :", "")
}

// GetRuntime 获取时长
func (s *FC2Scraper) GetRuntime() string {
	return "0"
}

// GetStudio 获取厂商
func (s *FC2Scraper) GetStudio() string {
	return util.FC2
}

// GetSeries 获取系列
func (s *FC2Scraper) GetSeries() string {
	return util.FC2
}

// GetTags 获取标签
func (s *FC2Scraper) GetTags() []string {
	// 组合地址
	uri := fmt.Sprintf("http://adult.contents.fc2.com/api/v4/article/%s/tag?", s.code)

	// 读取远程数据
	data, err := util.GetResult(uri, s.Proxy, nil)
	// 检查
	if err != nil {
		return nil
	}

	// 读取内容
	body, err := ioutil.ReadAll(bytes.NewReader(data))
	// 检查错误
	if err != nil {
		return nil
	}

	// json
	var tagsJSON fc2tags

	// 解析json
	err = json.Unmarshal(body, &tagsJSON)
	// 检查
	if err != nil {
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

// GetCover 获取图片
func (s *FC2Scraper) GetCover() string {
	// 获取图片
	fanart, _ := s.fc2clubRoot.Find(`.slides li:nth-child(1) img`).Attr("src")
	// 检查
	if fanart == "" {
		return ""
	}

	if fanart[0:2] == ".." {
		fanart = fanart[2:]
	}
	// 组合地址
	return fmt.Sprintf("https://fc2club.net%s", fanart)
}

// GetActors 获取演员
func (s *FC2Scraper) GetActors() map[string]string {
	return nil
}

// GetURI 获取页面地址
func (s *FC2Scraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *FC2Scraper) GetNumber() string {
	return s.number
}
