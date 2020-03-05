package scraper

// IScraper 刮削器接口
type IScraper interface {
	// Fetch 获取数据
	Fetch(code string) error

	// GetURI 获取来源页面地址
	GetURI() string

	// GetNumber 获取刮削番号
	GetNumber() string

	// GetTitle 获取影片名称
	GetTitle() string
	// GetIntro 获取影片简介
	GetIntro() string
	// GetDirector 获取影片导演
	GetDirector() string
	// GetRelease 获取发行时间
	GetRelease() string
	// GetRuntime 获取影片时长
	GetRuntime() string
	// GetStudio 获取影片厂商
	GetStudio() string
	// GetSeries 获取影片系列
	GetSeries() string
	// GetTags 获取标签列表
	GetTags() []string
	// GetCover 获取封面图片
	GetCover() string
	// GetActors 获取演员列表
	GetActors() map[string]string
}
