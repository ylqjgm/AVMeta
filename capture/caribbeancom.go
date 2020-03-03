package capture

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
)

// CaribBeanComCapture 加勒比刮削器
type CaribBeanComCapture struct {
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// Fetch 刮削
func (s *CaribBeanComCapture) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)

	// 组合地址
	uri := fmt.Sprintf("https://www.caribbeancom.com/moviepages/%s/index.html", code)

	// 打开远程连接
	data, err := GetResult(uri, s.Proxy, nil)
	// 检查
	if nil != err {
		return err
	}

	// 编码转换
	reader := transform.NewReader(bytes.NewReader(data), japanese.EUCJP.NewDecoder())

	// 获取根节点
	root, err := goquery.NewDocumentFromReader(reader)
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
func (s *CaribBeanComCapture) GetTitle() string {
	return s.root.Find(`h1[itemprop="name"]`).Text()
}

// GetIntro 获取简介
func (s *CaribBeanComCapture) GetIntro() string {
	// 获取简介
	intro, err := s.root.Find(`p[itemprop="description"]`).Html()
	// 检查
	if nil != err {
		return ""
	}

	return IntroFilter(intro)
}

// GetDirector 获取导演
func (s *CaribBeanComCapture) GetDirector() string {
	return ""
}

// GetRelease 发行时间
func (s *CaribBeanComCapture) GetRelease() string {
	return s.root.Find(`span[itemprop="uploadDate"]`).Text()
}

// GetRuntime 获取时长
func (s *CaribBeanComCapture) GetRuntime() string {
	// 获取数据
	strTime := strings.TrimSpace(s.root.Find(`span[itemprop="duration"]`).Text())

	// 是否正确获取
	if "" != strTime {
		// 搜索正则
		r, _ := regexp.Compile(`^(\d+):(\d+):(\d+)$`)
		// 搜索
		t := r.FindStringSubmatch(strTime)

		// 获取小时
		hour, err := strconv.Atoi(t[1])
		// 检查
		if nil != err {
			hour = 0
		}
		// 获取分钟
		minute, err := strconv.Atoi(t[2])
		// 检查
		if nil != err {
			minute = 0
		}

		return strconv.Itoa((hour * 60) + minute)
	}

	return "0"
}

// GetStudio 获取厂商
func (s *CaribBeanComCapture) GetStudio() string {
	return "カリビアンコム"
}

// GetSerise 获取系列
func (s *CaribBeanComCapture) GetSerise() string {
	return s.root.Find(`a[href*="/series/"]`).Text()
}

// GetTags 获取标签
func (s *CaribBeanComCapture) GetTags() []string {
	// 类别数组
	var tags []string
	// 循环获取
	s.root.Find(`a[itemprop="genre"]`).Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetFanart 获取图片
func (s *CaribBeanComCapture) GetFanart() string {
	return fmt.Sprintf("https://www.caribbeancom.com/moviepages/%s/images/l_l.jpg", s.number)
}

// GetActors 获取演员
func (s *CaribBeanComCapture) GetActors() map[string]string {
	// 演员列表
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`a[class="spec__tag"] span[itemprop="name"]`).Each(func(i int, item *goquery.Selection) {
		// 演员名称
		actors[strings.TrimSpace(item.Text())] = ""
	})

	return actors
}

// GetURI 获取页面地址
func (s *CaribBeanComCapture) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *CaribBeanComCapture) GetNumber() string {
	return s.number
}
