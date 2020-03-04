package capture

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// JavBusCapture javbus刮削器
type JavBusCapture struct {
	Site   string            // 免翻地址
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// Fetch 刮削
func (s *JavBusCapture) Fetch(code string) error {
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
func (s *JavBusCapture) detail() error {
	// 组合uri
	uri := fmt.Sprintf("%s/%s", CheckDomainPrefix(s.Site), s.number)
	// 获取节点
	root, err := GetRoot(uri, s.Proxy, nil)
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
func (s *JavBusCapture) GetTitle() string {
	return s.root.Find("h3").Text()
}

// GetIntro 获取简介
func (s *JavBusCapture) GetIntro() string {
	return GetDmmIntro(s.number, s.Proxy)
}

// GetDirector 获取导演
func (s *JavBusCapture) GetDirector() string {
	return s.root.Find(`a[href*="/director/"]`).Text()
}

// GetRelease 发行时间
func (s *JavBusCapture) GetRelease() string {
	return strings.ReplaceAll(s.root.Find(`p:contains("發行日期:")`).Text(), "發行日期: ", "")
}

// GetRuntime 获取时长
func (s *JavBusCapture) GetRuntime() string {
	return strings.ReplaceAll(strings.TrimRight(s.root.Find(`p:contains("長度:")`).Text(), "分鐘"), "長度: ", "")
}

// GetStudio 获取厂商
func (s *JavBusCapture) GetStudio() string {
	return s.root.Find(`a[href*="/studio/"]`).Text()
}

// GetSerise 获取系列
func (s *JavBusCapture) GetSerise() string {
	return s.root.Find(`a[href*="/series/"]`).Text()
}

// GetTags 获取标签
func (s *JavBusCapture) GetTags() []string {
	// 类别数组
	var tags []string
	// 循环获取
	s.root.Find(`span.genre a[href*="/genre/"]`).Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetFanart 获取图片
func (s *JavBusCapture) GetFanart() string {
	// 获取图片
	fanart, _ := s.root.Find(`a.bigImage img`).Attr("src")

	return fanart
}

// GetActors 获取演员
func (s *JavBusCapture) GetActors() map[string]string {
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
func (s *JavBusCapture) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *JavBusCapture) GetNumber() string {
	return s.number
}
