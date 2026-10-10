package main

import (
	"os"

	"github.com/amatsu-tech/amr/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
