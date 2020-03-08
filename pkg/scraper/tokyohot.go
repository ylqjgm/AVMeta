package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/PuerkitoBio/goquery"
)

// TokyoHotScraper tokyohot网站刮削器
type TokyoHotScraper struct {
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// NewTokyoHotScraper 返回一个被初始化的tokyohot刮削对象
//
// proxy 字符串参数，传入代理信息
func NewTokyoHotScraper(proxy string) *TokyoHotScraper {
	return &TokyoHotScraper{Proxy: proxy}
}

// Fetch 刮削
func (s *TokyoHotScraper) Fetch(code string) error {
	// 设置番号
	s.number = strings.ToLower(code)
	// 获取编号
	id, err := s.search()
	// 检查
	if id == "" || err != nil {
		return fmt.Errorf("404 Not Found")
	}

	// 组合地址
	uri := fmt.Sprintf("https://my.tokyo-hot.com%s?lang=zh-TW", id)
	// 打开链接
	root, err := util.GetRoot(uri, s.Proxy, nil)
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

// 搜索
func (s *TokyoHotScraper) search() (id string, err error) {
	// 组合地址
	uri := fmt.Sprintf("https://my.tokyo-hot.com/product/?q=%s&x=0&y=0&lang=zh-TW", s.number)
	// 获取节点
	root, err := util.GetRoot(uri, s.Proxy, nil)
	// 检查错误
	if err != nil {
		return
	}

	// 是否找到
	if -1 < root.Find(`ul.list > li:contains("沒有登入")`).Index() {
		err = fmt.Errorf("404 Not Found")
		return
	}

	// 获取结果
	root.Find(`ul.list li.detail a`).Each(func(i int, item *goquery.Selection) {
		// 获取番号
		number, _ := item.Find("img").Attr("title")
		// 转换大写
		number = strings.ToUpper(number)
		// 比较是否一致
		if !strings.EqualFold(strings.ToUpper(s.number), number) {
			return
		}
		// 获取地址链接
		id, _ = item.Attr("href")
	})

	// 检查是否获取到
	if id == "" {
		err = fmt.Errorf("404 Not Found")
		return
	}

	return id, err
}

// GetTitle 获取名称
func (s *TokyoHotScraper) GetTitle() string {
	return s.root.Find(`.pagetitle h2`).Text()
}

// GetIntro 获取简介
func (s *TokyoHotScraper) GetIntro() string {
	// 获取简介
	intro, err := s.root.Find(`div.sentence`).Html()
	// 检查错误
	if err != nil {
		return ""
	}

	return util.IntroFilter(intro)
}

// GetDirector 获取导演
func (s *TokyoHotScraper) GetDirector() string {
	return util.TOKYOHOT
}

// GetRelease 发行时间
func (s *TokyoHotScraper) GetRelease() string {
	return s.root.Find(`dt:contains("配信開始日")`).Next().Text()
}

// GetRuntime 获取时长
func (s *TokyoHotScraper) GetRuntime() string {
	// 获取时长
	strTime := strings.TrimSpace(s.root.Find(`dt:contains("収録時間"`).Next().Text())
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
func (s *TokyoHotScraper) GetStudio() string {
	return "東京熱"
}

// GetSeries 获取系列
func (s *TokyoHotScraper) GetSeries() string {
	return s.root.Find(`dt:contains("系列")`).Next().Find("a").Text()
}

// GetTags 获取标签
func (s *TokyoHotScraper) GetTags() []string {
	// 标签数组
	var tags []string
	// 循环获取
	s.root.Find(`dt:contains("Tag")`).Next().Find("a").Each(func(i int, item *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(item.Text()))
	})

	return tags
}

// GetCover 获取图片
func (s *TokyoHotScraper) GetCover() string {
	// 获取图片
	fanart, _ := s.root.Find(`.flowplayer video`).Attr("poster")

	return fanart
}

// GetActors 获取演员
func (s *TokyoHotScraper) GetActors() map[string]string {
	// 演员数组
	actors := make(map[string]string)

	// 循环获取
	s.root.Find(`dt:contains("出演者")`).Next().Find("a").Each(func(i int, item *goquery.Selection) {
		// 获取连接
		link, _ := item.Attr("href")
		// 组合地址
		uri := fmt.Sprintf("https://my.tokyo-hot.com%s", link)
		// 打开链接
		root, err := util.GetRoot(uri, s.Proxy, nil)
		// 检查错误
		if err != nil {
			return
		}

		// 获取演员图片
		img, _ := item.Find(`#profile img`).Attr("src")
		// 获取演员名字
		actors[strings.TrimSpace(root.Find(`.pagetitle h2`).Text())] = strings.TrimSpace(img)
	})

	return actors
}

// GetURI 获取页面地址
func (s *TokyoHotScraper) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *TokyoHotScraper) GetNumber() string {
	return s.number
}
