package main

import (
	"os"

	"github.com/ghas-projects/mrva-prep/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
