package cmd

import (
	"github.com/spf13/cobra"
)

func (e *Executor) initRoot() {
	rootCmd := &cobra.Command{
		Use:   "AVMeta",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}

	e.rootCmd = rootCmd
	e.initRootFlags()
}

func (e *Executor) initRootFlags() {
	e.rootCmd.Flags().BoolP("help", "h", false, "命令帮助")
	e.rootCmd.Flags().BoolP("version", "v", false, "程序版本查询")
}
