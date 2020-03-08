package scraper

import (
	"fmt"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

// JavDBScraper javdb网站刮削器
type JavDBScraper struct {
	Site   string            // 免翻地址
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// NewJavDBScraper 返回一个被初始化的javdb刮削对象
//
// site 字符串参数，传入免翻地址，
// proxy 字符串参数，传入代理信息
func NewJavDBScraper(site, proxy string) *JavDBScraper {
	return &JavDBScraper{Site: site, Proxy: proxy}
}

// Fetch 刮削
func (s *JavDBScraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 搜索
	id, err := s.search()
	// 检查错误
	if err != nil {
		return err
	}

	// 组合地址
	uri := fmt.Sprintf("%s%s", util.CheckDomainPrefix(s.Site), id)

	// 打开连接
	root, err := util.GetRoot(uri, s.Proxy, nil)
	// 检查错误
	if err != nil {
		return err
	}

	// 设置页面地址
	s.uri = uri
	// 设置根节点
	s.root = root

	return nil
}

// 搜索影片
func (s *JavDBScraper) search() (string, error) {
	// 组合地址
	uri := fmt.Sprintf("%s/search?q=%s&f=all", util.CheckDomainPrefix(s.Site), strings.ToUpper(s.number))

	// 打开地址
	root, err := util.GetRoot(uri, s.Proxy, nil)
	// 检查错误
	if err != nil {
		return "", err
	}

	// 查找是否获取到
	if -1 < root.Find(`.empty-message:contains("暫無內容")`).Index() {
		return "", fmt.Errorf("404 Not Found")
	}

	// 定义ID
	var id string

	// 循环检查番号
	root.Find(`div#videos .grid-item a`).Each(func(i int, item *goquery.Selection) {
		// 获取番号
		date := item.Find("div.uid").Text()
		// 大写并去除空白
		date = strings.ToUpper(strings.TrimSpace(date))
		// 检查番号是否完全正确
		if strings.EqualFold(strings.ToUpper(s.number), date) {
			// 获取href元素
			id, _ = item.Attr("href")
		}
	})

	// 清除空白
	id = strings.TrimSpace(id)

	// 是否获取到
	if id == "" {
		return "", fmt.Errorf("404 Not Found")
	}

	return id, nil
}

// GetTitle 获取名称
func (s *JavDBScraper) GetTitle() string {
	return strings.ReplaceAll(s.root.Find(`title`).Text(), "| JavDB 成人影片資料庫", "")
}

// GetIntro 获取简介
func (s *JavDBScraper) GetIntro() string {
	return GetDmmIntro(s.number, s.Proxy)
}

// GetDirector 获取导演
func (s *JavDBScraper) GetDirector() string {
	// 获取数据
	val := s.root.Find(`strong:contains("導演")`).Parent().NextFiltered(`span.value`).Text()
	// 检查
	if val == "" {
		val = s.root.Find(`strong:contains("導演")`).Parent().NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetRelease 发行时间
func (s *JavDBScraper) GetRelease() string {
	// 获取数据
	val := s.root.Find(`strong:contains("時間")`).Parent().NextFiltered(`span.value`).Text()
	// 检查
	if val == "" {
		val = s.root.Find(`strong:contains("時間")`).Parent().NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetRuntime 获取时长
func (s *JavDBScraper) GetRuntime() string {
	// 获取数据
	val := s.root.Find(`strong:contains("時長")`).Parent().NextFiltered(`span.value`).Text()
	// 检查
	if val == "" {
		val = s.root.Find(`strong:contains("時長")`).Parent().NextFiltered(`span.value`).Find("a").Text()
	}

	// 去除多余
	val = strings.TrimRight(val, "分鍾")

	return val
}

// GetStudio 获取厂商
func (s *JavDBScraper) GetStudio() string {
	// 获取数据
	val := s.root.Find(`strong:contains("片商")`).Parent().NextFiltered(`span.value`).Text()
	// 检查
	if val == "" {
		val = s.root.Find(`strong:contains("片商")`).Parent().NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetSeries 获取系列
func (s *JavDBScraper) GetSeries() string {
	// 获取数据
	val := s.root.Find(`strong:contains("系列")`).Parent().NextFiltered(`span.value`).Text()
	// 检查
	if val == "" {
		val = s.root.Find(`strong:contains("系列")`).Parent().NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetTags 获取标签
func (s *JavDBScraper) GetTags() []string {
	// 类别数组
	var tags []string
	// 循环获取
	s.root.Find(`strong:contains("类别")`).Parent().Parent().Find(".value a").Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetCover 获取图片
func (s *JavDBScraper) GetCover() string {
	// 获取图片
	fanart, _ := s.root.Find(`div.column-video-cover a img`).Attr("src")

	return fanart
}

// GetActors 获取演员
func (s *JavDBScraper) GetActors() map[string]string {
	// 演员列表
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`strong:contains("演員")`).Parent().Parent().Find(`.value a`).Each(func(i int, item *goquery.Selection) {
		// 演员名称
		actors[strings.TrimSpace(item.Text())] = ""
	})

	return actors
}

// GetURI 获取页面地址
func (s *JavDBScraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *JavDBScraper) GetNumber() string {
	return s.number
}
