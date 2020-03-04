package capture

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HeyzoCapture heyzo刮削器
type HeyzoCapture struct {
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	code   string            // 临时番号
	number string            // 最终番号
	root   *goquery.Document // 根节点
	json   *heyzoJSON        // json数据
}

// json结构
type heyzoJSON struct {
	Name            string `json:"name"`
	Image           string `json:"image"`
	DateCreated     string `json:"dateCreated"`
	Duration        string `json:"duration"`
	AggregateRating struct {
		Type        string `json:"@type"`
		RatingValue string `json:"ratingValue"`
		BestRating  string `json:"bestRating"`
		ReviewCount string `json:"reviewCount"`
	} `json:"aggregateRating"`
}

// Fetch 刮削
func (s *HeyzoCapture) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToUpper(code)
	// 番号正则
	r := regexp.MustCompile(`[0-9]{4}`)
	// 临时番号
	s.code = r.FindString(code)
	// 组合地址
	uri := fmt.Sprintf("https://www.heyzo.com/moviepages/%s/index.html", s.code)
	// 打开连接
	root, err := GetRoot(uri, s.Proxy, nil)
	// 检查
	if err != nil {
		return err
	}

	// 获取json节点
	data, err := root.Find(`script[type="application/ld+json"]`).Html()
	// 检查
	if err != nil {
		return err
	}
	// json对象
	js := &heyzoJSON{}
	// 转码
	data = strings.ReplaceAll(html.UnescapeString(data), "\n", "")
	// 转换为结构体
	err = json.Unmarshal([]byte(data), js)
	// 检查
	if err != nil {
		return err
	}

	// 设置页面地址
	s.uri = uri
	// 赋值根节点
	s.root = root
	// 赋值json
	s.json = js

	return nil
}

// GetTitle 获取名称
func (s *HeyzoCapture) GetTitle() string {
	return s.json.Name
}

// GetIntro 获取简介
func (s *HeyzoCapture) GetIntro() string {
	return IntroFilter(s.root.Find(`p[class="memo"]`).Text())
}

// GetDirector 获取导演
func (s *HeyzoCapture) GetDirector() string {
	return HEYZO
}

// GetRelease 发行时间
func (s *HeyzoCapture) GetRelease() string {
	return s.json.DateCreated
}

// GetRuntime 获取时长
func (s *HeyzoCapture) GetRuntime() string {
	// 获取时长
	duration := s.json.Duration
	// 时长搜索正则
	r := regexp.MustCompile(`^PT(\d+)H(\d+)M(\d+)S$`)
	// 搜索
	t := r.FindStringSubmatch(duration)
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

// GetStudio 获取厂商
func (s *HeyzoCapture) GetStudio() string {
	return "HEYZO"
}

// GetSerise 获取系列
func (s *HeyzoCapture) GetSerise() string {
	return s.root.Find(`.table-series a`).Text()
}

// GetTags 获取标签
func (s *HeyzoCapture) GetTags() []string {
	// 标签数组
	var tags []string
	// 循环获取
	s.root.Find(`.table-tag-keyword-big .tag-keyword-list li a`).Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetFanart 获取图片
func (s *HeyzoCapture) GetFanart() string {
	return "https:" + s.json.Image
}

// GetActors 获取演员
func (s *HeyzoCapture) GetActors() map[string]string {
	// 演员数组
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`.table-actor a span`).Each(func(i int, item *goquery.Selection) {
		// 获取演员名字
		actors[strings.TrimSpace(item.Text())] = ""
	})

	return actors
}

// GetURI 获取页面地址
func (s *HeyzoCapture) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *HeyzoCapture) GetNumber() string {
	return s.number
}
