package media

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/ylqjgm/AVMeta/pkg/scraper"
)

// 刮削对象
type captures struct {
	Name string
	S    scraper.IScraper
	R    *regexp.Regexp
}

// Pack 整理给定影片并返回 Media 结构体，
// 若整理失败则返回空对象及错误信息。
//
// file 字符串参数，传入要整理的文件路径，
// cfg ConfigStruct结构体，传入程序配置信息。
func Pack(file string, cfg *util.ConfigStruct) (*Media, error) {
	if cfg.Media.Library == "vsmeta" {
		return packVSMeta(file, cfg)
	}

	return packNfo(file, cfg)
}

// 整理给定影片为 nfo 并返回 Media 结构体，
// 若整理失败则返回空对象及错误信息。
//
// file 字符串参数，传入要整理的文件路径，
// cfg ConfigStruct结构体，传入程序配置信息。
func packNfo(file string, cfg *util.ConfigStruct) (*Media, error) {
	// 获取采集数据
	m, err := capture(file, cfg)
	// 检查
	if err != nil {
		return nil, err
	}

	// 转换为XML
	buff, err := mediaToXML(m)
	// 检查
	if err != nil {
		return nil, err
	}

	// 写入nfo
	err = util.WriteFile(fmt.Sprintf("%s/%s.nfo", m.DirPath, m.Number), buff)
	// 检查
	if err != nil {
		return nil, err
	}

	// 获取视频后缀
	ext := path.Ext(file)
	// 移动视频文件
	err = util.MoveFile(file, fmt.Sprintf("%s/%s%s", m.DirPath, m.Number, ext))

	return m, err
}

// 整理给定影片为 vsmeta 并返回 VSMeta 结构体，
// 若整理失败则返回空对象及错误信息。
//
// file 字符串参数，传入要整理的文件路径，
// cfg ConfigStruct结构体，传入程序配置信息。
func packVSMeta(file string, cfg *util.ConfigStruct) (*Media, error) {
	// 获取整理数据
	m, err := capture(file, cfg)
	// 检查
	if err != nil {
		return nil, err
	}

	// 实例化VSMeta
	vs := NewVSMeta()
	// 解析为 vsmeta
	vs.ParseVSMeta(m)
	// 写入封面
	vs.writePoster(fmt.Sprintf("%s/poster.jpg", m.DirPath))
	// 写入背景
	vs.writeFanart(fmt.Sprintf("%s/fanart.jpg", m.DirPath))

	// 获取视频后缀
	ext := path.Ext(file)

	// 写入vsmeta
	err = util.WriteFile(fmt.Sprintf("%s/%s%s.vsmeta", m.DirPath, m.Number, ext), vs.B.Bytes())
	// 检查
	if err != nil {
		return nil, err
	}

	// 删除封面
	os.Remove(fmt.Sprintf("%s/poster.jpg", m.DirPath))
	// 删除背景
	os.Remove(fmt.Sprintf("%s/fanart.jpg", m.DirPath))

	// 移动视频文件
	err = util.MoveFile(file, fmt.Sprintf("%s/%s%s", m.DirPath, m.Number, ext))

	return m, err
}

// 整理影片并返回 Media 对象
//
// file 字符串参数，传入要整理的文件路径，
// cfg ConfigStruct结构体，传入程序配置信息。
func capture(file string, cfg *util.ConfigStruct) (*Media, error) {
	// 搜索番号并获得刮削对象
	m, err := search(file, cfg)
	// 检查
	if err != nil {
		return nil, err
	}

	// 是否有图片
	if m.Cover == "" {
		return nil, fmt.Errorf("找不到封面")
	}

	// 获取准确目录
	dirPath := util.GetNumberPath(m.ConvertMap(), cfg)
	// 创建目录
	err = os.MkdirAll(dirPath, os.ModePerm)
	// 检查
	if err != nil {
		return nil, err
	}
	// 赋值保存路径
	m.DirPath = dirPath

	// 获取图片后缀
	ext := path.Ext(m.Cover)
	// 下载图片
	err = util.SavePhoto(m.Cover, fmt.Sprintf("%s/fanart.jpg", m.DirPath), cfg.Base.Proxy, !strings.EqualFold(strings.ToLower(ext), ".jpg"))
	// 检查
	if err != nil {
		return nil, err
	}

	// 裁剪图片
	err = util.PosterCover(fmt.Sprintf("%s/fanart.jpg", m.DirPath), fmt.Sprintf("%s/poster.jpg", m.DirPath), cfg)
	// 检查
	if err != nil {
		return nil, err
	}

	// 设定图片
	m.FanArt = fmt.Sprintf("fanart.jpg")
	m.Poster = fmt.Sprintf("poster.jpg")
	m.Thumb = fmt.Sprintf("poster.jpg")

	return m, nil
}

// 番号搜索
func search(file string, cfg *util.ConfigStruct) (*Media, error) {
	// 定义变量
	var err error

	// 提取番号
	code := util.GetCode(file, cfg.Path.Filter)

	// 定义一个拥有正则匹配的刮削对象数组
	sr := []captures{
		{
			Name: "CaribBeanCom",
			S:    scraper.NewCaribBeanComScraper(cfg.Base.Proxy),
			R:    regexp.MustCompile(`^\d{6}-\d{3}$`),
		},
		{
			Name: "TokyoHot",
			S:    scraper.NewTokyoHotScraper(cfg.Base.Proxy),
			R:    regexp.MustCompile(`(^red-\d{3}|n\d{4})`),
		},
		{
			Name: "Heyzo",
			S:    scraper.NewHeyzoScraper(cfg.Base.Proxy),
			R:    regexp.MustCompile(`^heyzo-[0-9]{4}`),
		},
		{
			Name: "Heydouga",
			S:    scraper.NewHeydougaScraper(cfg.Base.Proxy),
			R:    regexp.MustCompile(`([0-9]{4}).+?([0-9]{3,4})$`),
		},
		{
			Name: "FC2",
			S:    scraper.NewFC2Scraper(cfg.Base.Proxy),
			R:    regexp.MustCompile(`^fc2-[0-9]{6,7}`),
		},
		{
			Name: "Siro",
			S:    scraper.NewSiroScraper(cfg.Base.Proxy),
			R:    regexp.MustCompile(`^(siro|abp|[0-9]{3,4}[a-zA-Z]{2,5})-[0-9]{3,4}`),
		},
		{
			Name: "DMM",
			S:    scraper.NewDMMScraper(cfg.Base.Proxy),
			R:    regexp.MustCompile(`[a-zA-Z]{2,5}[-|\s\S][0-9]{3,4}`),
		},
	}
	// 定义一个没有正则匹配的刮削对象数组
	ss := []captures{
		{
			Name: "JavDB",
			S:    scraper.NewJavDBScraper(cfg.Site.JavDB, cfg.Base.Proxy),
			R:    nil,
		},
		{
			Name: "JavBus",
			S:    scraper.NewJavBusScraper(cfg.Site.JavBus, cfg.Base.Proxy),
			R:    nil,
		},
	}

	// 转换番号为小写
	code = strings.ToLower(code)
	// 定义一个刮削对象
	var s scraper.IScraper

	// 查找正则匹配
	for _, scr := range sr {
		// 刮削赋值
		s = scr.S
		// 检查是否匹配
		if scr.R.MatchString(code) {
			// 刮削
			err = s.Fetch(code)
			break
		}
	}

	// 检查错误
	if err != nil || s == nil {
		// 尝试刮削
		for _, sc := range ss {
			// 刮削赋值
			s = sc.S
			// 刮削
			if err = s.Fetch(code); err == nil {
				break
			}
		}
	}

	// 再次检测
	if err != nil || s == nil {
		return nil, err
	}

	// 刮削并获取nfo对象
	return ParseMedia(s)
}

// 转换为xml
func mediaToXML(m *Media) ([]byte, error) {
	// 转换
	x, err := xml.MarshalIndent(m, "", "  ")
	// 检查
	if err != nil {
		return nil, err
	}

	// 转码为[]byte
	x = []byte(xml.Header + string(x))

	return x, nil
}
