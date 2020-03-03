package capture

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// JavDBCapture javdb刮削器
type JavDBCapture struct {
	Site   string            // 免翻地址
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// Fetch 刮削
func (s *JavDBCapture) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 搜索
	id, err := s.search()
	// 检查错误
	if nil != err {
		return err
	}

	// 组合地址
	uri := fmt.Sprintf("%s%s", CheckDomainPrefix(s.Site), id)

	// 打开连接
	root, err := GetRoot(uri, s.Proxy, nil)
	// 检查错误
	if nil != err {
		return err
	}

	// 设置页面地址
	s.uri = uri
	// 设置根节点
	s.root = root

	return nil
}

// 搜索影片
func (s *JavDBCapture) search() (string, error) {
	// 组合地址
	uri := fmt.Sprintf("%s/search?q=%s&f=all", CheckDomainPrefix(s.Site), strings.ToUpper(s.number))

	// 打开地址
	root, err := GetRoot(uri, s.Proxy, nil)
	// 检查错误
	if nil != err {
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
		if strings.ToUpper(s.number) == date {
			// 获取href元素
			id, _ = item.Attr("href")
		}
	})

	// 清除空白
	id = strings.TrimSpace(id)

	// 是否获取到
	if "" == id {
		return "", fmt.Errorf("404 Not Found")
	}

	return id, nil
}

// GetTitle 获取名称
func (s *JavDBCapture) GetTitle() string {
	return strings.ReplaceAll(s.root.Find(`title`).Text(), "| JavDB 成人影片資料庫", "")
}

// GetIntro 获取简介
func (s *JavDBCapture) GetIntro() string {
	return GetDmmIntro(s.number, s.Proxy)
}

// GetDirector 获取导演
func (s *JavDBCapture) GetDirector() string {
	// 获取数据
	val := s.root.Find(`strong:contains("導演")`).Parent().NextFiltered(`span.value`).Text()
	// 检查
	if "" == val {
		val = s.root.Find(`strong:contains("導演")`).Parent().NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetRelease 发行时间
func (s *JavDBCapture) GetRelease() string {
	// 获取数据
	val := s.root.Find(`strong:contains("時間")`).Parent().NextFiltered(`span.value`).Text()
	// 检查
	if "" == val {
		val = s.root.Find(`strong:contains("時間")`).Parent().NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetRuntime 获取时长
func (s *JavDBCapture) GetRuntime() string {
	// 获取数据
	val := s.root.Find(`strong:contains("時長")`).Parent().NextFiltered(`span.value`).Text()
	// 检查
	if "" == val {
		val = s.root.Find(`strong:contains("時長")`).Parent().NextFiltered(`span.value`).Find("a").Text()
	}

	// 去除多余
	val = strings.TrimRight(val, "分鍾")

	return val
}

// GetStudio 获取厂商
func (s *JavDBCapture) GetStudio() string {
	// 获取数据
	val := s.root.Find(`strong:contains("片商")`).Parent().NextFiltered(`span.value`).Text()
	// 检查
	if "" == val {
		val = s.root.Find(`strong:contains("片商")`).Parent().NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetSerise 获取系列
func (s *JavDBCapture) GetSerise() string {
	// 获取数据
	val := s.root.Find(`strong:contains("系列")`).Parent().NextFiltered(`span.value`).Text()
	// 检查
	if "" == val {
		val = s.root.Find(`strong:contains("系列")`).Parent().NextFiltered(`span.value`).Find("a").Text()
	}

	return val
}

// GetTags 获取标签
func (s *JavDBCapture) GetTags() []string {
	// 类别数组
	var tags []string
	// 循环获取
	s.root.Find(`strong:contains("类别")`).Parent().Parent().Find(".value a").Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetFanart 获取图片
func (s *JavDBCapture) GetFanart() string {
	// 获取图片
	fanart, _ := s.root.Find(`div.column-video-cover a img`).Attr("src")

	return fanart
}

// GetActors 获取演员
func (s *JavDBCapture) GetActors() map[string]string {
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
func (s *JavDBCapture) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *JavDBCapture) GetNumber() string {
	return s.number
}
