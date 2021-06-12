package scraper

import (
	"fmt"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

// JavBusScraper javbus网站刮削器
type JavBusScraper struct {
	Site   string            // 免翻地址
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// NewJavBusScraper 返回一个被初始化的javbus刮削对象
//
// site 字符串参数，传入免翻地址，
// proxy 字符串参数，传入代理信息
func NewJavBusScraper(site, proxy string) *JavBusScraper {
	return &JavBusScraper{Site: site, Proxy: proxy}
}

// Fetch 刮削
func (s *JavBusScraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 获取信息
	err := s.detail()
	// 检查错误
	if err != nil {
		// 设置番号
		s.number = strings.ReplaceAll(s.number, "-", "_")
		// 使用 _ 方式
		err = s.detail()
		// 检查错误
		if err != nil {
			// 设置番号
			s.number = strings.ReplaceAll(strings.ReplaceAll(s.number, "-", ""), "_", "")
			// 去除符号
			err = s.detail()
			// 检查错误
			if err != nil {
				return fmt.Errorf("404 Not Found")
			}
		}
	}

	return nil
}

// 获取获取
func (s *JavBusScraper) detail() error {
	// 组合uri
	uri := fmt.Sprintf("%s/%s", util.CheckDomainPrefix(s.Site), s.number)
	// 获取节点
	root, err := util.GetRoot(uri, s.Proxy, nil)
	// 检查错误
	if err != nil {
		return err
	}

	// 查找是否获取到
	if -1 == root.Find(`h3`).Index() {
		return fmt.Errorf("404 Not Found")
	}

	// 设置页面地址
	s.uri = uri
	// 设置根节点
	s.root = root

	return nil
}

// GetTitle 获取名称
func (s *JavBusScraper) GetTitle() string {
	return s.root.Find("h3").Text()
}

// GetIntro 获取简介
func (s *JavBusScraper) GetIntro() string {
	return GetDmmIntro(s.number, s.Proxy)
}

// GetDirector 获取导演
func (s *JavBusScraper) GetDirector() string {
	return s.root.Find(`a[href*="/director/"]`).Text()
}

// GetRelease 发行时间
func (s *JavBusScraper) GetRelease() string {
	return strings.ReplaceAll(s.root.Find(`p:contains("發行日期:")`).Text(), "發行日期: ", "")
}

// GetRuntime 获取时长
func (s *JavBusScraper) GetRuntime() string {
	return strings.ReplaceAll(strings.TrimRight(s.root.Find(`p:contains("長度:")`).Text(), "分鐘"), "長度: ", "")
}

// GetStudio 获取厂商
func (s *JavBusScraper) GetStudio() string {
	return s.root.Find(`a[href*="/studio/"]`).Text()
}

// GetSeries 获取系列
func (s *JavBusScraper) GetSeries() string {
	return s.root.Find(`a[href*="/series/"]`).Text()
}

// GetTags 获取标签
func (s *JavBusScraper) GetTags() []string {
	// 类别数组
	var tags []string
	// 循环获取
	s.root.Find(`span.genre a[href*="/genre/"]`).Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetCover 获取图片
func (s *JavBusScraper) GetCover() string {
	// 获取图片
	fanart, _ := s.root.Find(`a.bigImage img`).Attr("src")

	if find := strings.Contains(fanart, s.Site); !find {
		fanart = s.Site + fanart
	}

	return fanart
}

// GetActors 获取演员
func (s *JavBusScraper) GetActors() map[string]string {
	// 演员数组
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`div.star-box li > a`).Each(func(i int, item *goquery.Selection) {
		// 获取演员图片
		img, _ := item.Find(`img`).Attr("src")
		// 获取演员名字
		name, _ := item.Find("img").Attr("title")

		// 加入列表
		actors[strings.TrimSpace(name)] = strings.TrimSpace(img)
	})

	return actors
}

// GetURI 获取页面地址
func (s *JavBusScraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *JavBusScraper) GetNumber() string {
	return s.number
}
