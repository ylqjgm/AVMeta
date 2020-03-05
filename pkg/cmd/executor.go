package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Executor 命令对象
type Executor struct {
	rootCmd *cobra.Command

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
	e.initRoot()
	e.setTemplate()
	e.initVersion()

	return e
}

// Execute 执行命令
func (e *Executor) Execute() error {
	return e.rootCmd.Execute()
}
