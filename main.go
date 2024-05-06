package main

import (
	"os"

	"github.com/catatsuy/purl/internal/cli"
	"golang.org/x/term"
)

func main() {
	c := cli.NewCLI(os.Stdout, os.Stderr, os.Stdin, term.IsTerminal(int(os.Stdin.Fd())), term.IsTerminal(int(os.Stdout.Fd())))
	os.Exit(c.Run(os.Args))
}
