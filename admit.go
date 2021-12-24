package main

import (
	"os"

	"github.com/qqbuby/kube-admit/pkg/cmd"
)

func main() {
	root := cmd.NewAdmitCommand()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
