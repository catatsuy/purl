package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"golang.org/x/term"
)

const (
	ExitCodeOK             = 0
	ExitCodeParseFlagError = 1
	ExitCodeFail           = 1
)

type CLI struct {
	outStream, errStream io.Writer
	inputStream          io.Reader
}

func NewCLI(errStream io.Writer, inputStream io.Reader) *CLI {
	return &CLI{errStream: errStream, inputStream: inputStream}
}

func (c *CLI) Run(args []string) int {
	flags := flag.NewFlagSet("notify_slack", flag.ContinueOnError)
	flags.SetOutput(c.errStream)

	var replaceExpr string
	var inplaceEdit bool

	flags.BoolVar(&inplaceEdit, "i", false, "overwrite the file inplace")
	flags.StringVar(&replaceExpr, "replace", "", `Replacement expression, e.g., "@search@replace@"`)

	err := flags.Parse(args[1:])
	if err != nil {
		return ExitCodeParseFlagError
	}

	filePath := ""
	if flags.NArg() > 0 {
		filePath = flags.Arg(0)
	} else if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintf(c.errStream, "No input file specified\n")
		return ExitCodeFail
	}

	if inplaceEdit && filePath == "" {
		fmt.Fprintf(c.errStream, "Cannot use -i option with stdin\n")
		return ExitCodeFail
	}

	if len(replaceExpr) < 3 {
		fmt.Fprintf(c.errStream, "Invalid replace expression format. Use \"@search@replace@\"\n")
		return ExitCodeFail
	}

	delimiter := string(replaceExpr[0])
	parts := regexp.MustCompile(regexp.QuoteMeta(delimiter)).Split(replaceExpr[1:], -1)
	if len(parts) < 2 {
		fmt.Fprintf(c.errStream, "Invalid replace expression format. Use \"@search@replace@\"\n")
		return ExitCodeFail
	}
	searchPattern, replacement := parts[0], parts[1]

	var scanner *bufio.Scanner

	if filePath == "" {
		scanner = bufio.NewScanner(c.inputStream)
	} else {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to open file: %s\n", err)
			return ExitCodeFail
		}
		defer file.Close()

		scanner = bufio.NewScanner(file)
	}

	var tmpFile *os.File

	if inplaceEdit {
		tmpFile, err = os.CreateTemp("", "purl")
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to create temp file: %s\n", err)
			return ExitCodeFail
		}
		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		c.outStream, err = os.Create(tmpFile.Name())
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to open file for writing: %s\n", err)
			return ExitCodeFail
		}
	} else {
		c.outStream = os.Stdout
	}

	re, err := regexp.Compile(searchPattern)
	if err != nil {
		fmt.Fprintf(c.errStream, "Invalid regex pattern: %s\n", err)
		return ExitCodeFail
	}

	for scanner.Scan() {
		line := scanner.Text()
		modifiedLine := re.ReplaceAllString(line, replacement)
		fmt.Fprintf(c.outStream, modifiedLine+"\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(c.errStream, "Error reading file: %s\n", err)
		return ExitCodeFail
	}

	if inplaceEdit {
		if err := os.Rename(tmpFile.Name(), filePath); err != nil {
			fmt.Fprintf(c.errStream, "Failed to overwrite the original file: %s\n", err)
			return ExitCodeFail
		}
	}

	return ExitCodeOK
}
