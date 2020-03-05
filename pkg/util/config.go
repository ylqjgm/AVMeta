package util

import (
	"github.com/spf13/viper"
)

// BaseStruct 基础配置结构
type BaseStruct struct {
	Proxy string // 代理地址
}

// PathStruct 路径配置结构
type PathStruct struct {
	Success   string // 成功存储目录
	Fail      string // 失败存储目录
	Directory string // 影片存储路径格式
	Filter    string // 文件名过滤规则
}

// MediaStruct 媒体库配置结构
type MediaStruct struct {
	Library   string // 媒体库类型
	URL       string // Emby访问地址
	API       string // Emby API Key
	SecretID  string // 腾讯云 SecretId
	SecretKey string // 腾讯云 SecretKey
}

// SiteStruct 免费地址配置结构
type SiteStruct struct {
	JavBus string // javbus免翻地址
	JavDB  string // javdb免翻地址
}

// ConfigStruct 配置结构
type ConfigStruct struct {
	Base  BaseStruct  // 基础配置
	Path  PathStruct  // 路径配置
	Media MediaStruct // 媒体库配置
	Site  SiteStruct  // 免翻地址配置
}

// GetConfig 读取配置信息
func GetConfig() (*ConfigStruct, error) {
	// 配置名称
	viper.SetConfigName("config")
	// 配置类型
	viper.SetConfigType("yaml")
	// 添加当前执行路径为配置路径
	viper.AddConfigPath(".")
	// 读取配置信息
	err := viper.ReadInConfig()
	// 读取配置
	if err != nil {
		// 如果文件不存在
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return WriteConfig()
		}

		// 直接返回错误信息
		return nil, err
	}

	// 定义配置变量
	var config ConfigStruct

	// 反序列
	err = viper.Unmarshal(&config)

	return &config, err
}

// WriteConfig 写入配置文件
func WriteConfig() (*ConfigStruct, error) {
	// 配置名称
	viper.SetConfigName("config")
	// 配置类型
	viper.SetConfigType("yaml")
	// 添加当前执行路径为配置路径
	viper.AddConfigPath(".")

	// 默认配置
	cfg := &ConfigStruct{
		Base: BaseStruct{
			Proxy: "",
		},
		Path: PathStruct{
			Success:   "success",
			Fail:      "fail",
			Directory: "{actor}/{year}/{number}",
			Filter:    "-hd||hd-||_hd||hd_||[||]||【||】||asfur||~||-full||_full||3xplanet||monv||云中飘荡||@||tyhg999.com||xxxxxxxx||-fhd||_fhd||thz.la",
		},
		Media: MediaStruct{
			Library:   "emby",
			URL:       "",
			API:       "",
			SecretID:  "",
			SecretKey: "",
		},
		Site: SiteStruct{
			JavBus: "https://www.javbus.com/",
			JavDB:  "https://javdb4.com/",
		},
	}

	// 设置数据
	viper.Set("base", cfg.Base)
	viper.Set("path", cfg.Path)
	viper.Set("media", cfg.Media)
	viper.Set("site", cfg.Site)

	return cfg, viper.SafeWriteConfig()
}
