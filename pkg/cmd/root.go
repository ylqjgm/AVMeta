package cmd

import (
	"fmt"
	"log"
	"path"

	"github.com/schollz/progressbar/v2"

	"github.com/ylqjgm/AVMeta/pkg/media"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/spf13/cobra"
)

func (e *Executor) initRoot() {
	e.rootCmd = &cobra.Command{
		Use:   "AVMeta",
		Short: "一款使用 Golang 编写的跨平台 AV 元数据刮削器",
		Long: `
AVMeta 是一款使用 Golang 编写的跨平台 AV 元数据刮削器
使用 AVMeta, 您可自动将 AV 电影进行归类整理
并生成对应媒体库元数据文件`,
		Run: rootRunFunc,
	}
}

func (e *Executor) setTemplate() {
	// 重设使用显示模板
	e.rootCmd.SetUsageTemplate(`使用:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if .HasExample}}

示例:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

命令:
  help        命令执行帮助
  config      生成配置文件
  version     显示程序版本{{end}}{{if .HasAvailableSubCommands}}

使用 "{{.CommandPath}} help [command]" 可获取更多命令帮助.{{end}}

`)
}

// root命令执行函数
func rootRunFunc(cmd *cobra.Command, args []string) {
	// 获取配置信息
	cfg, err := util.GetConfig()
	// 检查错误
	if err != nil {
		log.Fatalln(err)
	}

	// 获取当前执行路径
	curDir := util.GetRunPath()

	// 列当前目录
	files, err := util.WalkDir(curDir, cfg.Path.Success, cfg.Path.Fail)
	// 检测错误
	if err != nil {
		log.Fatalln(err)
	}

	// 获取总量
	count := len(files)
	// 定义两个计数变量
	fail := 0
	success := 0
	// 输出总量
	fmt.Printf("\n共探索到 %d 个视频文件, 开始刮削整理...\n\n", count)

	// 定义进度条
	bar := progressbar.NewOptions(count,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetDescription(fmt.Sprintf("[blue][%d/%d][reset] ...", 0, count)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[cyan]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[cyan][[reset]",
			BarEnd:        "[cyan]][reset]",
		}),
	)

	// 循环视频文件列表
	for _, file := range files {
		// 输出文字
		bar.Describe(fmt.Sprintf("[[light_red]%d[reset]/[green]%d[reset]] 「[cyan]%s[reset]」...", fail, success, path.Base(file)))
		// 进度条自增
		_ = bar.Add(util.ONE)
		// 刮削整理
		m, err := media.Pack(file, cfg)
		// 检查
		if err != nil {
			fmt.Println(err)
			// 错误
			fail++
			// 恢复文件
			util.FailFile(file, cfg)
			continue
		}
		// 成功
		success++
		// 输出文字
		bar.Describe(fmt.Sprintf("[[light_red]%d[reset]/[green]%d[reset]] 「[cyan]%s[reset]」...", fail, success, m.Number))
	}

	// 进度条结束
	_ = bar.Finish()
	// 输出一条空行
	fmt.Println("")
}
