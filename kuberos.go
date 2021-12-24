package main

import (
	"os"

	"github.com/qqbuby/kuberos/pkg/cmd"
)

func main() {
	root := cmd.NewKuberosCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
