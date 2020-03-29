package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

// SiroScraper siro网站刮削器
type SiroScraper struct {
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// NewSiroScraper 返回一个被初始化的siro刮削对象
//
// proxy 字符串参数，传入代理信息
func NewSiroScraper(proxy string) *SiroScraper {
	return &SiroScraper{Proxy: proxy}
}

// Fetch 刮削
func (s *SiroScraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 定义Cookies
	var cookies []*http.Cookie
	// 加入Cookie
	cookies = append(cookies, &http.Cookie{
		Name:    "adc",
		Value:   "1",
		Path:    "/",
		Domain:  "mgstage.com",
		Expires: time.Now().Add(1 * time.Hour),
	})
	// 组合地址
	uri := fmt.Sprintf("https://www.mgstage.com/product/product_detail/%s/", s.number)
	// 打开链接
	root, err := util.GetRoot(uri, s.Proxy, cookies)
	// 检查
	if err != nil {
		return err
	}

	// 设置页面地址
	s.uri = uri
	// 设置根节点
	s.root = root

	return nil
}

// GetTitle 获取名称
func (s *SiroScraper) GetTitle() string {
	return s.root.Find(`h1.tag`).Text()
}

// GetIntro 获取简介
func (s *SiroScraper) GetIntro() string {
	return util.IntroFilter(s.root.Find(`#introduction p.introduction`).Text())
}

// GetDirector 获取导演
func (s *SiroScraper) GetDirector() string {
	return ""
}

// GetRelease 发行时间
func (s *SiroScraper) GetRelease() string {
	return s.root.Find(`th:contains("配信開始日")`).NextFiltered("td").Text()
}

// GetRuntime 获取时长
func (s *SiroScraper) GetRuntime() string {
	return strings.TrimRight(s.root.Find(`th:contains("収録時間")`).NextFiltered("td").Text(), "min")
}

// GetStudio 获取厂商
func (s *SiroScraper) GetStudio() string {
	return s.root.Find(`th:contains("メーカー")`).NextFiltered("td").Text()
}

// GetSeries 获取系列
func (s *SiroScraper) GetSeries() string {
	return s.root.Find(`th:contains("シリーズ")`).NextFiltered("td").Text()
}

// GetTags 获取标签
func (s *SiroScraper) GetTags() []string {
	// 标签数组
	var tags []string
	// 循环获取
	s.root.Find(`th:contains("ジャンル")`).NextFiltered("td").Find("a").Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetCover 获取图片
func (s *SiroScraper) GetCover() string {
	// 获取图片
	fanart, _ := s.root.Find(`#EnlargeImage`).Attr("href")

	return fanart
}

// GetActors 获取演员
func (s *SiroScraper) GetActors() map[string]string {
	// 演员数组
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`th:contains("出演")`).NextFiltered("td").Find("a").Each(func(i int, item *goquery.Selection) {
		// 演员名字
		actors[strings.TrimSpace(item.Text())] = ""
	})

	// 是否获取到
	if len(actors) == 0 {
		// 重新获取
		name := s.root.Find(`th:contains("出演")`).NextFiltered("td").Text()
		// 获取演员名字
		actors[strings.TrimSpace(name)] = ""
	}

	return actors
}

// GetURI 获取页面地址
func (s *SiroScraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *SiroScraper) GetNumber() string {
	return s.number
}
