package capture

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HeydougaCapture heydouga刮削器
type HeydougaCapture struct {
	Proxy  string            // 代理配置
	uri    string            // 页面地址
	data   string            // 页面数据
	code   string            // 临时番号
	code1  string            // 前面部分
	code2  string            // 后面部分
	number string            // 最终番号
	root   *goquery.Document // 根节点
}

// json结构
type heydougaJSON struct {
	Tag []heydougaTag `json:"tag"`
}

// 标签结构
type heydougaTag struct {
	TagName string `json:"tag_name"`
}

// Fetch 刮削
func (s *HeydougaCapture) Fetch(code string) error {
	// 设置临时番号
	s.code = code
	// 转换大写
	code = strings.ToUpper(code)
	// 番号正则
	r, _ := regexp.Compile(`([0-9]{4}).+?([0-9]{3,4})`)
	// 临时番号
	code = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(r.FindString(code), "PPV", ""), "HEYDOUGA", ""))
	// 检查是否为空
	if "" == code {
		return fmt.Errorf("找不到番号")
	}

	// 配置番号
	code = r.FindString(code)
	// 番号分割
	cs := strings.Split(code, "-")
	// 检查是否有两个
	if 2 != len(cs) {
		return fmt.Errorf("找不到番号")
	}

	// 设置番号前后缀
	s.code1 = cs[0]
	s.code2 = cs[1]
	// 组合地址
	uri := fmt.Sprintf("https://www.heydouga.com/moviepages/%s/%s/index.html", s.code1, s.code2)
	// 打开连接
	data, status, err := MakeRequest("GET", uri, s.Proxy, nil, nil, nil)
	// 检查
	if nil != err || http.StatusNotFound == status {
		// 设置番号前后缀
		s.code1 = cs[0]
		s.code2 = "ppv-" + cs[1]
		// 重新组合地址
		uri = fmt.Sprintf("https://www.heydouga.com/moviepages/%s/%s/index.html", s.code1, s.code2)
		// 打开链接
		data, status, err = MakeRequest("GET", uri, s.Proxy, nil, nil, nil)
		// 检查
		if nil != err || http.StatusNotFound == status {
			return err
		}
	}

	// 获取根节点
	root, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	// 检查
	if nil != err {
		return err
	}

	// 设置番号
	s.number = fmt.Sprintf("Heydouga %s-PPV%s", s.code1, s.code2)
	// 设置页面地址
	s.uri = uri
	// 赋值根节点
	s.root = root
	// 设置页面数据
	s.data = string(data)

	return nil
}

// GetTitle 获取名称
func (s *HeydougaCapture) GetTitle() string {
	return s.root.Find(`div#title-bg h1`).Text()
}

// GetIntro 获取简介
func (s *HeydougaCapture) GetIntro() string {
	return IntroFilter(s.root.Find(`div[class="movie-description"] p`).Text())
}

// GetDirector 获取导演
func (s *HeydougaCapture) GetDirector() string {
	return s.root.Find(`span:contains("提供元")`).Next().Find(`a[href*="/listpages/provider"]`).Text()
}

// GetRelease 发行时间
func (s *HeydougaCapture) GetRelease() string {
	return s.root.Find(`span:contains("配信日")`).Next().Text()
}

// GetRuntime 获取时长
func (s *HeydougaCapture) GetRuntime() string {
	return strings.ReplaceAll(s.root.Find(`span:contains("動画再生時間")`).Next().Text(), "分", "")
}

// GetStudio 获取厂商
func (s *HeydougaCapture) GetStudio() string {
	return "Hey動画"
}

// GetSerise 获取系列
func (s *HeydougaCapture) GetSerise() string {
	return "Hey動画 PPV"
}

// GetTags 获取标签
func (s *HeydougaCapture) GetTags() []string {
	// movie_seq正则表达式
	r, _ := regexp.Compile(`movie_seq:([0-9]+)`)
	mm := r.FindString(s.data)
	fmt.Println(mm)
	// 搜索movie_seq
	m := strings.TrimSpace(strings.ReplaceAll(r.FindString(s.data), "movie_seq:", ""))
	// 检查是否获取到
	if "" == m {
		return nil
	}

	// 组合路径
	uri := fmt.Sprintf("https://www.heydouga.com/get_movie_tag_all_utf8/?movie_seq=%s", m)
	// 获取数据
	data, err := GetResult(uri, s.Proxy, nil)
	// 检查错误
	if nil != err {
		return nil
	}

	// json对象
	js := &heydougaJSON{}
	// 转换为结构体
	err = json.Unmarshal(data, js)
	// 检查
	if nil != err {
		return nil
	}

	// 定义标签数组
	var tags []string
	// 循环标签
	for _, tag := range js.Tag {
		// 加入标签
		tags = append(tags, strings.TrimSpace(tag.TagName))
	}

	return tags
}

// GetFanart 获取图片
func (s *HeydougaCapture) GetFanart() string {
	return fmt.Sprintf("https://image01-www.heydouga.com/contents/%s/%s/player_thumb.jpg", s.code1, s.code2)
}

// GetActors 获取演员
func (s *HeydougaCapture) GetActors() map[string]string {
	// 演员map
	actors := make(map[string]string)
	// 定义一个临时演员数组
	var tmpActors []string

	// 循环获取
	s.root.Find(`span:contains("主演")`).Next().Find(`a`).Each(func(i int, item *goquery.Selection) {
		// 获取演员信息
		act := strings.TrimSpace(item.Text())
		// 检查
		if "" == act {
			return
		}
		// 分割数据
		acts1 := strings.Split(act, "、")
		acts2 := strings.Split(act, " ")
		// 循环加入数组
		for _, a := range acts1 {
			tmpActors = append(tmpActors, strings.TrimSpace(strings.ReplaceAll(a, "素人", "")))
		}
		for _, a := range acts2 {
			tmpActors = append(tmpActors, strings.TrimSpace(strings.ReplaceAll(a, "素人", "")))
		}
	})

	// 循环加入map
	for _, actor := range tmpActors {
		actors[actor] = ""
	}

	return actors
}

// GetURI 获取页面地址
func (s *HeydougaCapture) GetURI() string {
	return s.uri
}

// GetNumber 获取番号
func (s *HeydougaCapture) GetNumber() string {
	return s.number
}
