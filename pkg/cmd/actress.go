package cmd

import (
	"github.com/ylqjgm/AVMeta/pkg/logs"
	"strings"

	"github.com/ylqjgm/AVMeta/pkg/actress"

	"github.com/spf13/cobra"
)

// actress命令
func (e *Executor) initActress() {
	actressCmd := &cobra.Command{
		Use: "actress",
		Long: `
自动从各网站提取女优头像并上传至 Emby 服务器中`,
		Example: `  AVMeta actress
  AVMeta actress --site javbus
  AVMeta actress down --site javbus
  AVMeta actress put`,
		Run: e.actressRunFunc,
	}

	// 添加参数
	actressCmd.Flags().StringVar(&site, "site", "", "采集站点: javbus, javdb")
	e.rootCmd.AddCommand(actressCmd)
}

// 头像执行命令
func (e *Executor) actressRunFunc(cmd *cobra.Command, args []string) {
	// 初始化日志
	logs.Log("actress")

	// 定义参数变量
	var arg string
	down := false

	// 检测参数
	if len(args) > 1 {
		// 输出帮助
		_ = cmd.Help()
		return
	} else if len(args) > 0 {
		// 获取参数
		arg = args[0]

		// 检查参数
		if !strings.EqualFold(arg, "down") {
			_ = cmd.Help()
			return
		}

		down = arg == "down"
	}

	// 是否配置了 Emby 数据
	if e.cfg.Media.URL == "" || e.cfg.Media.API == "" {
		logs.Fatal("Emby 访问地址或 API Key 未配置, 请配置后重试")
	}

	// 是否为入库
	if !down {
		// 初始化对象
		actor := actress.NewActress()
		// 入库头像
		logs.Info("开始入库本地女优头像...")
		// 调用入库
		_ = actor.Put()

		return
	}

	// 如果设置站点
	if site != "" {
		// 转大写
		site = strings.ToUpper(site)

		// 检查传入参数正确性
		if len(site) > 0 && site != actress.JAVDB && site != actress.JAVBUS {
			logs.Fatal("--site 参数仅支持 javbus, javdb 两个选项, 留空则全部采集.")
		}
	}

	// 如果是下载
	if down {
		// 仅javBUS
		if site == actress.JAVBUS {
			// 下载javBUS
			fetchJavBUS()
			return
		}

		// 仅javDB
		if site == actress.JAVDB {
			// 下载javDB
			fetchJavDB()
			return
		}

		fetchJavBUS()
		fetchJavDB()
	}
}

// 下载javBUS
func fetchJavBUS() {
	// 初始化对象
	actor := actress.NewActress()
	// 下载javbus有码
	logs.Info("开始下载 JavBus 有码女优头像...")
	_ = actor.Fetch("JAVBUS", 1, true)
	// 下载javbus无码
	logs.Info("开始下载 JavBus 无码女优头像...")
	_ = actor.Fetch("JAVBUS", 1, false)
}

// 下载javDB
func fetchJavDB() {
	// 初始化对象
	actor := actress.NewActress()
	// 下载javdb有码
	logs.Info("开始下载 JavDB 有码女优头像...")
	_ = actor.Fetch("JAVDB", 1, true)
	// 下载javdb无码
	logs.Info("开始下载 JavDB 无码女优头像...")
	_ = actor.Fetch("JAVDB", 1, false)
}
