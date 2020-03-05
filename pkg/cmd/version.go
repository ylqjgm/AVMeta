package cmd

import (
	"github.com/spf13/cobra"
)

func (e *Executor) initVersion() {
	e.rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Long: `
执行本命令打印当前您所使用的 AVMeta 版本信息`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Version:\t%s\n", e.version)
			cmd.Printf("Git commit:\t%s\n", e.commit)
			cmd.Printf("Built:\t\t%s\n", e.built)
			cmd.Printf("Go Version:\t%s\n", e.goVersion)
			cmd.Printf("Platform:\t%s\n", e.platform)
		},
	})
}
