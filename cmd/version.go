package cmd

import (
	"github.com/spf13/cobra"
)

func (e *Executor) initVersion() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Version:\t%s\n", e.version)
			cmd.Printf("Git commit:\t%s\n", e.commit)
			cmd.Printf("Built:\t\t%s\n", e.built)
			cmd.Printf("Go Version:\t%s\n", e.goVersion)
			cmd.Printf("Platform:\t%s\n", e.platform)
		},
	}

	e.rootCmd.AddCommand(versionCmd)
}
