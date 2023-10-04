package main

import (
	"fmt"
	"os"

	"github.com/replicate/yolo/pkg/cli"
)

func main() {
	cmd, err := cli.NewRootCommand()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
