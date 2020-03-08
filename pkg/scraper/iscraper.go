/*
Package scraper 网站数据刮削器.

scraper通过IScraper接口，对具体网站进行数据刮削，
执行刮削操作后，将页面转换为树结构，并存储在刮削对象中，
当需要任何数据时，执行对应方法，从树结构中查找对应信息。
*/
package scraper

// IScraper 刮削器接口
type IScraper interface {
	// Fetch 执行刮削，并返回刮削结果
	//
	// code 字符串参数，传入番号信息
	Fetch(code string) error

	// GetURI 获取刮削的页面地址
	GetURI() string

	// GetNumber 获取最终的正确番号信息
	GetNumber() string

	// GetTitle 从刮削结果中获取影片标题
	GetTitle() string
	// GetIntro 从刮削结果中获取影片简介
	GetIntro() string
	// GetDirector 从刮削结果中获取影片导演
	GetDirector() string
	// GetRelease 从刮削结果中获取发行时间
	GetRelease() string
	// GetRuntime 从刮削结果中获取影片时长
	GetRuntime() string
	// GetStudio 从刮削结果中获取影片厂商
	GetStudio() string
	// GetSeries 从刮削结果中获取影片系列
	GetSeries() string
	// GetTags 从刮削结果中获取影片标签
	GetTags() []string
	// GetCover 从刮削结果中获取背景图片
	GetCover() string
	// GetActors 从刮削结果中获取影片演员
	GetActors() map[string]string
}
