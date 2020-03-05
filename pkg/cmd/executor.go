package cmd

import (
	"fmt"
	"log"
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

// NewExecutor 创建命令对象
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
	e.initVersion()

	return e
}

// Execute 执行命令
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
			log.Fatalln(err)
		}
	}

	e.cfg = cfg
}
