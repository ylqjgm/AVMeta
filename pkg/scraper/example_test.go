package scraper_test

import (
	"fmt"

	"github.com/ylqjgm/AVMeta/pkg/scraper"
)

// 数据刮削示例
func Example() {
	// 定义刮削结构
	var s scraper.IScraper

	// 以dmm作为演示，实例化一个刮削对象
	s = scraper.NewDMMScraper("")
	// 执行刮削操作
	err := s.Fetch("bf-592")
	if err != nil {
		panic(err)
	}

	// 获取标题
	fmt.Println(s.GetTitle())

	// Output:
	// BF-592 篠田ゆう4時間 SUPERBEST
}
