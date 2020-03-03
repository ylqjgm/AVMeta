package capture

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// SiroCapture siro刮削器
type SiroCapture struct {
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// Fetch 刮削
func (s *SiroCapture) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 定义Cookies
	var cookies []http.Cookie
	// 加入Cookie
	cookies = append(cookies, http.Cookie{
		Name:    "adc",
		Value:   "1",
		Path:    "/",
		Domain:  "mgstage.com",
		Expires: time.Now().Add(1 * time.Hour),
	})
	// 组合地址
	uri := fmt.Sprintf("https://www.mgstage.com/product/product_detail/%s/", s.number)
	// 打开链接
	root, err := GetRoot(uri, s.Proxy, cookies)
	// 检查
	if nil != err {
		return err
	}

	// 设置页面地址
	s.uri = uri
	// 设置根节点
	s.root = root

	return nil
}

// GetTitle 获取名称
func (s *SiroCapture) GetTitle() string {
	return s.root.Find(`h1.tag`).Text()
}

// GetIntro 获取简介
func (s *SiroCapture) GetIntro() string {
	return IntroFilter(s.root.Find(`#introduction p.introduction`).Text())
}

// GetDirector 获取导演
func (s *SiroCapture) GetDirector() string {
	return ""
}

// GetRelease 发行时间
func (s *SiroCapture) GetRelease() string {
	return s.root.Filter(`th:contains("配信開始日")`).Next().Text()
}

// GetRuntime 获取时长
func (s *SiroCapture) GetRuntime() string {
	return strings.TrimRight(s.root.Find(`th:contains("収録時間")`).Next().Text(), "min")
}

// GetStudio 获取厂商
func (s *SiroCapture) GetStudio() string {
	return s.root.Find(`th:contains("メーカー")`).Next().Text()
}

// GetSerise 获取系列
func (s *SiroCapture) GetSerise() string {
	return s.root.Find(`th:contains("シリーズ")`).Next().Text()
}

// GetTags 获取标签
func (s *SiroCapture) GetTags() []string {
	// 标签数组
	var tags []string
	// 循环获取
	s.root.Find(`th:contains("ジャンル")`).Next().Find("a").Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetFanart 获取图片
func (s *SiroCapture) GetFanart() string {
	// 获取图片
	fanart, _ := s.root.Find(`#EnlargeImage`).Attr("href")

	return fanart
}

// GetActors 获取演员
func (s *SiroCapture) GetActors() map[string]string {
	// 演员数组
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`th:contains("出演")`).Next().Find("a").Each(func(i int, item *goquery.Selection) {
		// 演员名字
		actors[strings.TrimSpace(item.Text())] = ""
	})

	// 是否获取到
	if 0 >= len(actors) {
		// 重新获取
		name := s.root.Find(`th:contains("出演")`).Next().Text()
		// 获取演员名字
		actors[strings.TrimSpace(name)] = ""
	}

	return actors
}

// GetURI 获取页面地址
func (s *SiroCapture) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *SiroCapture) GetNumber() string {
	return s.number
}
