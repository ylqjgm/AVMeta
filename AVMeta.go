package main

import (
	"github.com/ylqjgm/AVMeta/pkg/cmd"
)

var (
	version = "master"
	commit  = "?"
	built   = ""
)

func main() {
	e := cmd.NewExecutor(version, commit, built)

	if err := e.Execute(); err != nil {
		panic(err)
	}
}
