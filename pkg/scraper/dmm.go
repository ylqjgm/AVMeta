package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

// DMMScraper dmm网站刮削器
type DMMScraper struct {
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	code   string            // 临时番号
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// NewDMMScraper 返回一个被初始化的dmm刮削对象
//
// proxy 字符串参数，传入代理信息
func NewDMMScraper(proxy string) *DMMScraper {
	return &DMMScraper{Proxy: proxy}
}

// Fetch 刮削
func (s *DMMScraper) Fetch(code string) error {
	// 大写
	code = strings.ToUpper(code)
	// 设置番号
	s.number = code
	// 查询所用番号
	code = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(code, "_", ""), "-", ""))
	// 临时番号
	s.code = code

	// 获取根节点
	err := s.getRoot()

	return err
}

// GetDmmIntro 从dmm网站中获取影片简介。
//
// code 字符串参数，传入番号，
// proxy 字符串参数，传入代理信息
func GetDmmIntro(code, proxy string) string {
	// 实例化对象
	s := NewDMMScraper(proxy)
	// 获取数据
	err := s.Fetch(code)
	// 检查
	if err != nil {
		return ""
	}

	return s.GetIntro()
}

// 获取根节点
func (s *DMMScraper) getRoot() error {
	// 组合地址列表
	uris := []string{
		"https://www.dmm.co.jp/digital/videoa/-/detail/=/cid=%s",
		"https://www.dmm.co.jp/mono/dvd/-/detail/=/cid=%s",
		"https://www.dmm.co.jp/digital/anime/-/detail/=/cid=%s",
		"https://www.dmm.co.jp/mono/anime/-/detail/=/cid=%s",
	}

	// 定义Cookies
	var cookies []*http.Cookie
	// 加入年龄确认Cookie
	cookies = append(cookies, &http.Cookie{
		Name:  "age_check_done",
		Value: "1",
	})

	var err error
	var root *goquery.Document

	// 循环
	for _, uri := range uris {
		// 打开连接
		root, err = util.GetRoot(fmt.Sprintf(uri, s.code), s.Proxy, cookies)
		// 检查
		if err == nil {
			// 设置页面地址
			s.uri = uri
			// 设置根节点
			s.root = root

			// 判断是否返回了地域限制
			foreignError := root.Find(`.foreignError__desc`).Text()
			if foreignError != "" {
				return fmt.Errorf(foreignError)
			}
			return nil
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// GetTitle 获取名称
func (s *DMMScraper) GetTitle() string {
	return s.root.Find(`h1#title`).Text()
}

// GetIntro 获取简介
func (s *DMMScraper) GetIntro() string {
	return util.IntroFilter(s.root.Find(`tr td div.mg-b20.lh4 p.mg-b20`).Text())
}

// GetDirector 获取导演
func (s *DMMScraper) GetDirector() string {
	// 获取导演
	director := s.root.Find(`td:contains("監督：")`).Next().Find("a").Text()
	// 如果没有
	if director == "" {
		director = s.root.Find(`td:contains("監督：")`).Next().Text()
	}

	return director
}

// GetRelease 发行时间
func (s *DMMScraper) GetRelease() string {
	// 获取发行时间
	release := s.root.Find(`td:contains("発売日：")`).Next().Find("a").Text()
	// 没获取到
	if release == "" {
		release = s.root.Find(`td:contains("発売日：")`).Next().Text()
	}

	// 替换
	release = strings.ReplaceAll(release, "/", "-")

	return release
}

// GetRuntime 获取时长
func (s *DMMScraper) GetRuntime() string {
	return strings.TrimRight(s.root.Find(`td:contains("収録時間：")`).Next().Text(), "分")
}

// GetStudio 获取厂商
func (s *DMMScraper) GetStudio() string {
	// 获取厂商
	studio := s.root.Find(`td:contains("メーカー：")`).Next().Find("a").Text()
	// 是否获取到
	if studio == "" {
		studio = s.root.Find(`td:contains("メーカー：")`).Next().Text()
	}

	return studio
}

// GetSeries 获取系列
func (s *DMMScraper) GetSeries() string {
	// 获取系列
	set := s.root.Find(`td:contains("シリーズ：")`).Next().Find("a").Text()
	// 是否获取到
	if set == "" {
		set = s.root.Find(`td:contains("シリーズ：")`).Next().Text()
	}

	return set
}

// GetTags 获取标签
func (s *DMMScraper) GetTags() []string {
	// 标签数组
	var tags []string
	// 循环获取
	s.root.Find(`td:contains("ジャンル：")`).Next().Find("a").Each(func(i int, item *goquery.Selection) {
		// 加入数组
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetCover 获取图片
func (s *DMMScraper) GetCover() string {
	// 获取图片
	fanart, _ := s.root.Find(`#` + s.code).Attr("href")

	if fanart == "" {
		s.root.Find(`td:contains("品番：")`).Next().Each(func(i int, item *goquery.Selection) {
			// 获取演员名字
			number := strings.TrimSpace(item.Text())
			if number != "" {
				cover, _ := s.root.Find(`#` + number).Attr("href")
				fanart = cover
			}
		})
	}

	return fanart
}

// GetActors 获取演员
func (s *DMMScraper) GetActors() map[string]string {
	// 演员数组
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`td:contains("出演者：")`).Next().Find("span a").Each(func(i int, item *goquery.Selection) {
		// 获取演员名字
		actors[strings.TrimSpace(item.Text())] = ""
	})

	return actors
}

// GetURI 获取页面地址
func (s *DMMScraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *DMMScraper) GetNumber() string {
	return s.number
}
