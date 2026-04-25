package main

import (
	"os"

	"github.com/overiss/go-boilerplater/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args))
}
