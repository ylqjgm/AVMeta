/*
Package cmd 命令行操作包。

AVMeta程序所有操作命令皆由此包定义，使用 cobra 第三方包编写。
*/
package cmd

import (
	"fmt"
	"github.com/ylqjgm/AVMeta/pkg/logs"
	"runtime"

	"github.com/ylqjgm/AVMeta/pkg/util"

	"github.com/spf13/cobra"
)

// 采集站点变量
var site string

// Executor 命令对象
type Executor struct {
	rootCmd *cobra.Command
	cfg     *util.ConfigStruct

	version   string
	commit    string
	built     string
	goVersion string
	platform  string
}

// NewExecutor 返回一个被初始化的命令对象。
//
// version 字符串参数，传入当前程序版本，
// commit 字符串参数，传入最后提交的 git commit，
// built 字符串参数，传入程序编译时间。
func NewExecutor(version, commit, built string) *Executor {
	e := &Executor{
		version:   version,
		commit:    commit,
		built:     built,
		goVersion: runtime.Version(),
		platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
	e.initConfig()
	e.initRoot()
	e.setTemplate()
	e.initConfigFile()
	e.initActress()
	e.initNfo()
	e.initVersion()

	return e
}

// Execute 执行根命令。
func (e *Executor) Execute() error {
	return e.rootCmd.Execute()
}

// 初始化配置
func (e *Executor) initConfig() {
	// 获取配置
	cfg, err := util.GetConfig()
	// 检查
	if err != nil {
		// 初始化配置
		cfg, err = util.WriteConfig()
		// 检查
		if err != nil {
			logs.FatalError(err)
		}
	}

	// 配置信息
	e.cfg = cfg
}
