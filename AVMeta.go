package main

import (
	"fmt"
	"os"

	"github.com/ylqjgm/AVMeta/cmd"
)

var (
	version = "master"
	commit  = "?"
	built   = ""
)

func main() {
	e := cmd.NewExecutor(version, commit, built)

	if err := e.Execute(); err != nil {
		fmt.Printf("failed executing command with error %v\n", err)
		os.Exit(0)
	}
}
