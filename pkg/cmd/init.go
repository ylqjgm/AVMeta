package cmd

import "github.com/spf13/cobra"

func (e *Executor) initConfigFile() {
	e.rootCmd.AddCommand(&cobra.Command{
		Use: "init",
		Long: `
在当前目录下生成 config.yaml 配置文件`,
		Run: func(cmd *cobra.Command, args []string) {
			e.initConfig()
		},
	})
}
