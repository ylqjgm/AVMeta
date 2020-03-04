package capture

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// AVSoxCapture avsox刮削器
type AVSoxCapture struct {
	Site   string            // 免翻地址
	Proxy  string            // 代理地址
	number string            // 番号
	uri    string            // 来源地址
	root   *goquery.Document // 页面节点
}

// Fetch 获取数据
func (s *AVSoxCapture) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 搜索番号
	uri, err := s.search()
	// 检查错误
	if err != nil {
		return err
	}

	// 打开连接并获取节点
	root, err := GetRoot(uri, s.Proxy, nil)
	// 检查错误
	if err != nil {
		return err
	}

	// 赋值来源地址
	s.uri = uri
	// 赋值节点
	s.root = root

	return nil
}

// GetURI 获取来源页面地址
func (s *AVSoxCapture) GetURI() string { return s.uri }

// GetNumber 获取刮削番号
func (s *AVSoxCapture) GetNumber() string { return s.number }

// GetTitle 获取影片名称
func (s *AVSoxCapture) GetTitle() string {
	return s.root.Find(`h3`).Text()
}

// GetIntro 获取影片简介
func (s *AVSoxCapture) GetIntro() string {
	// 通过dmm获取
	return GetDmmIntro(s.number, s.Proxy)
}

// GetDirector 获取影片导演
func (s *AVSoxCapture) GetDirector() string {
	return s.root.Find(`a[href*="/director/"]`).Text()
}

// GetRelease 获取发行时间
func (s *AVSoxCapture) GetRelease() string {
	return strings.TrimLeft(s.root.Find(`span:contains("发行时间:")`).Parent().Text(), "发行时间: ")
}

// GetRuntime 获取影片时长
func (s *AVSoxCapture) GetRuntime() string {
	return strings.TrimLeft(strings.TrimRight(s.root.Find(`span:contains("长度:")`).Parent().Text(), "分钟"), "长度: ")
}

// GetStudio 获取影片厂商
func (s *AVSoxCapture) GetStudio() string {
	return s.root.Find(`a[href*="/studio/"]`).Text()
}

// GetSerise 获取影片系列
func (s *AVSoxCapture) GetSerise() string {
	return s.root.Find(`a[href*="/series/"]`).Text()
}

// GetTags 获取标签列表
func (s *AVSoxCapture) GetTags() []string {
	// 标签数组
	var tags []string
	// 循环获取
	s.root.Find(`span.genre a`).Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetCover 获取封面图片
func (s *AVSoxCapture) GetCover() string {
	// 获取图片
	cover, _ := s.root.Find(`a.bigImage img`).Attr("src")

	return cover
}

// GetActors 获取演员列表
func (s *AVSoxCapture) GetActors() map[string]string {
	// 演员数组
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`a.avatar-box`).Each(func(i int, item *goquery.Selection) {
		// 获取演员图片
		img, _ := item.Find(`img`).Attr("src")
		// 获取演员名字
		actors[strings.TrimSpace(item.Find(`span`).Text())] = strings.TrimSpace(img)
	})

	return actors
}

// 搜索番号
func (s *AVSoxCapture) search() (string, error) {
	// 获取页面id
	id, err := s.getID()
	// 检查错误
	if id == "" || err != nil {
		// 将 - 转换为 _
		s.number = strings.ReplaceAll(s.number, "-", "_")
		// 重新获取
		id, err = s.getID()
		// 检查错误
		if id == "" || err != nil {
			// 去除番号
			s.number = strings.ReplaceAll(s.number, "_", "")
			// 再次获取
			id, err = s.getID()
			// 检查错误
			if id == "" || err != nil {
				return "", fmt.Errorf("404 Not Found")
			}
		}
	}

	return id, err
}

// 获取页面编号
func (s *AVSoxCapture) getID() (string, error) {
	// 组合地址
	uri := fmt.Sprintf("%s/cn/search/%s", CheckDomainPrefix(s.Site), s.number)
	// 获取节点
	root, err := GetRoot(uri, s.Proxy, nil)
	// 检查错误
	if err != nil {
		return "", err
	}

	// 检查是否有数据
	if -1 < root.Find(`h4:contains("搜寻没有结果")`).Index() {
		return "", fmt.Errorf("404 Not Fond")
	}

	// 编号变量
	var id string

	// 循环检查
	root.Find(`div#waterfall .item`).Each(func(i int, item *goquery.Selection) {
		// 获取番号
		tmpNumber := item.Find("date").Eq(0).Text()
		// 转大写并去除空白
		tmpNumber = strings.ToUpper(strings.TrimSpace(tmpNumber))
		// 检查是否为传入番号
		if tmpNumber == s.number {
			// 获取编号
			id, _ = item.Find("a").Attr("href")
		}
	})

	// 去除空白
	id = strings.TrimSpace(id)

	// 检查
	if id == "" {
		return "", fmt.Errorf("404 Not Found")
	}

	return id, nil
}
