package media_test

import (
	"fmt"

	"github.com/ylqjgm/AVMeta/pkg/media"
	"github.com/ylqjgm/AVMeta/pkg/scraper"
	"github.com/ylqjgm/AVMeta/pkg/util"
)

// 转换刮削对象示例
func ExampleParseNfo() {
	// 初始化一个刮削对象
	var s scraper.IScraper

	// 以DMM演示
	s = scraper.NewDMMScraper("")
	// 刮削
	err := s.Fetch("BF-592")
	if err != nil {
		panic(err)
	}
	// 转换为 nfo 结构
	m, err := media.ParseNfo(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(m.Title)

	// Output:
	// BF-592 篠田ゆう4時間 SUPERBEST
}

// 转换路径参数示例
func ExampleMedia_ConvertMap() {
	// 初始化一个刮削对象
	var s scraper.IScraper

	// 以DMM演示
	s = scraper.NewDMMScraper("")
	// 刮削
	err := s.Fetch("BF-592")
	if err != nil {
		panic(err)
	}
	// 转换为 nfo 结构
	m, err := media.ParseNfo(s)
	if err != nil {
		panic(err)
	}

	// 转换部分对象
	filter := m.ConvertMap()

	fmt.Println(filter["{studio}"])

	// Output:
	// BeFree
}

// 文件整理示例
func ExamplePack() {
	// 首先获取配置信息
	cfg, err := util.GetConfig()
	// 本例以当前目录下存在一个 bf-592.mp4 的文件为前提
	m, err := media.Pack("./bf-592.mp4", cfg)
	if err != nil {
		panic(err)
	}

	// 到此，若未出错，将在当前目录下按照配置信息创建了对应目录，
	// 且已将 bf-592.mp4 文件移动到目录中，并下载了背景图片，
	// 剪切了封面图片，生成了 BF-592.nfo 的元数据文件。

	fmt.Println(m.Title)

	// Output:
	// BF-592 篠田ゆう4時間 SUPERBEST
}
