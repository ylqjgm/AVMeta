package scraper

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
)

// CaribBeanComScraper 加勒比网站刮削器
type CaribBeanComScraper struct {
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// NewCaribBeanComScraper 返回一个被初始化的加勒比刮削对象
//
// proxy 字符串参数，传入代理信息
func NewCaribBeanComScraper(proxy string) *CaribBeanComScraper {
	return &CaribBeanComScraper{Proxy: proxy}
}

// Fetch 刮削
func (s *CaribBeanComScraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)

	// 组合地址
	uri := fmt.Sprintf("https://www.caribbeancom.com/moviepages/%s/index.html", code)

	// 打开远程连接
	data, err := util.GetResult(uri, s.Proxy, nil)
	// 检查
	if err != nil {
		return err
	}

	// 编码转换
	reader := transform.NewReader(bytes.NewReader(data), japanese.EUCJP.NewDecoder())

	// 获取根节点
	root, err := goquery.NewDocumentFromReader(reader)
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

// GetTitle 获取标题
func (s *CaribBeanComScraper) GetTitle() string {
	return s.root.Find(`h1[itemprop="name"]`).Text()
}

// GetIntro 获取简介
func (s *CaribBeanComScraper) GetIntro() string {
	// 获取简介
	intro, err := s.root.Find(`p[itemprop="description"]`).Html()
	// 检查
	if err != nil {
		return ""
	}

	return util.IntroFilter(intro)
}

// GetDirector 获取导演
func (s *CaribBeanComScraper) GetDirector() string {
	return ""
}

// GetRelease 发行时间
func (s *CaribBeanComScraper) GetRelease() string {
	return s.root.Find(`span[itemprop="uploadDate"]`).Text()
}

// GetRuntime 影片时长
func (s *CaribBeanComScraper) GetRuntime() string {
	// 获取数据
	strTime := strings.TrimSpace(s.root.Find(`span[itemprop="duration"]`).Text())

	// 是否正确获取
	if strTime != "" {
		// 搜索正则
		r := regexp.MustCompile(`^(\d+):(\d+):(\d+)$`)
		// 搜索
		t := r.FindStringSubmatch(strTime)

		// 获取小时
		hour, err := strconv.Atoi(t[1])
		// 检查
		if err != nil {
			hour = 0
		}
		// 获取分钟
		minute, err := strconv.Atoi(t[2])
		// 检查
		if err != nil {
			minute = 0
		}

		return strconv.Itoa((hour * 60) + minute)
	}

	return "0"
}

// GetStudio 获取厂商
func (s *CaribBeanComScraper) GetStudio() string {
	return "カリビアンコム"
}

// GetSeries 影片系列
func (s *CaribBeanComScraper) GetSeries() string {
	return s.root.Find(`a[href*="/series/"]`).Text()
}

// GetTags 获取标签
func (s *CaribBeanComScraper) GetTags() []string {
	// 类别数组
	var tags []string
	// 循环获取
	s.root.Find(`a[itemprop="genre"]`).Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetCover 背景图片
func (s *CaribBeanComScraper) GetCover() string {
	return fmt.Sprintf("https://www.caribbeancom.com/moviepages/%s/images/l_l.jpg", s.number)
}

// GetActors 获取演员
func (s *CaribBeanComScraper) GetActors() map[string]string {
	// 演员列表
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`a[class="spec__tag"] span[itemprop="name"]`).Each(func(i int, item *goquery.Selection) {
		// 演员名称
		actors[strings.TrimSpace(item.Text())] = ""
	})

	return actors
}

// GetURI 页面地址
func (s *CaribBeanComScraper) GetURI() string {
	return s.uri
}

// GetNumber 正确番号
func (s *CaribBeanComScraper) GetNumber() string {
	return s.number
}
