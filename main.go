package main

import (
	"os"

	"github.com/catatsuy/purl/cli"
)

func main() {
	c := cli.NewCLI(os.Stdout, os.Stderr, os.Stdin)
	os.Exit(c.Run(os.Args))
}
