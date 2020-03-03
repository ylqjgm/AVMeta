package capture

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// DMMCapture dmm刮削器
type DMMCapture struct {
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	code   string            // 临时番号
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// Fetch 刮削
func (s *DMMCapture) Fetch(code string) error {
	// 大写
	code = strings.ToUpper(code)
	// 设置番号
	s.number = code
	// 查询所用番号
	code = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(code, "_", ""), "-", ""))
	// 临时番号
	s.code = code

	// 组合地址
	uri := fmt.Sprintf("https://www.dmm.co.jp/digital/videoa/-/detail/=/cid=%s", code)
	// 打开连接
	root, err := GetRoot(uri, s.Proxy, nil)
	// 检查
	if nil != err {
		// 重新组合地址
		uri = fmt.Sprintf("https://www.dmm.co.jp/mono/dvd/-/detail/=/cid=%s", code)
		// 打开连接
		root, err = GetRoot(uri, s.Proxy, nil)
		// 再次检查错误
		if nil != err {
			return err
		}
	}

	// 设置页面地址
	s.uri = uri
	// 设置根节点
	s.root = root

	return nil
}

// GetDmmIntro 直接获取dmm的简介
func GetDmmIntro(code, proxy string) string {
	// 查询所用番号
	code = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(code, "_", ""), "-", ""))

	// 组合地址
	uri := fmt.Sprintf("https://www.dmm.co.jp/digital/videoa/-/detail/=/cid=%s", code)
	// 打开连接
	root, err := GetRoot(uri, proxy, nil)
	// 检查
	if nil != err {
		// 重新组合地址
		uri = fmt.Sprintf("https://www.dmm.co.jp/mono/dvd/-/detail/=/cid=%s", code)
		// 打开连接
		root, err = GetRoot(uri, proxy, nil)
		// 再次检查错误
		if nil != err {
			return ""
		}
	}

	return IntroFilter(root.Find(`tr td div.mg-b20.lh4 p.mg-b20`).Text())
}

// GetTitle 获取名称
func (s *DMMCapture) GetTitle() string {
	return s.root.Find(`h1#title`).Text()
}

// GetIntro 获取简介
func (s *DMMCapture) GetIntro() string {
	return IntroFilter(s.root.Find(`tr td div.mg-b20.lh4 p.mg-b20`).Text())
}

// GetDirector 获取导演
func (s *DMMCapture) GetDirector() string {
	// 获取导演
	director := s.root.Find(`td:contains("監督：")`).Next().Find("a").Text()
	// 如果没有
	if "" == director {
		director = s.root.Find(`td:contains("監督：")`).Next().Text()
	}

	return director
}

// GetRelease 发行时间
func (s *DMMCapture) GetRelease() string {
	// 获取发行时间
	release := s.root.Find(`td:contains("発売日：")`).Next().Find("a").Text()
	// 没获取到
	if "" == release {
		release = s.root.Find(`td:contains("発売日：")`).Next().Text()
	}

	// 替换
	release = strings.ReplaceAll(release, "/", "-")

	return release
}

// GetRuntime 获取时长
func (s *DMMCapture) GetRuntime() string {
	return strings.TrimRight(s.root.Find(`td:contains("収録時間")`).Next().Text(), "分")
}

// GetStudio 获取厂商
func (s *DMMCapture) GetStudio() string {
	// 获取厂商
	studio := s.root.Find(`td:contains("メーカー")`).Next().Find("a").Text()
	// 是否获取到
	if "" == studio {
		studio = s.root.Find(`td:contains("メーカー")`).Next().Text()
	}

	return studio
}

// GetSerise 获取系列
func (s *DMMCapture) GetSerise() string {
	// 获取系列
	set := s.root.Find(`td:contains("シリーズ：")`).Next().Find("a").Text()
	// 是否获取到
	if "" == set {
		set = s.root.Find(`td:contains("シリーズ：")`).Next().Text()
	}

	return set
}

// GetTags 获取标签
func (s *DMMCapture) GetTags() []string {
	// 标签数组
	var tags []string
	// 循环获取
	s.root.Find(`td:contains("ジャンル：")`).Next().Find("a").Each(func(i int, item *goquery.Selection) {
		// 加入数组
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetFanart 获取图片
func (s *DMMCapture) GetFanart() string {
	// 获取图片
	fanart, _ := s.root.Find(`#` + s.code).Attr("href")

	return fanart
}

// GetActors 获取演员
func (s *DMMCapture) GetActors() map[string]string {
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
func (s *DMMCapture) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *DMMCapture) GetNumber() string {
	return s.number
}
