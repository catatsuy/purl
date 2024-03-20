package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
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

func NewCLI(outStream, errStream io.Writer, inputStream io.Reader) *CLI {
	return &CLI{outStream: outStream, errStream: errStream, inputStream: inputStream}
}

func (c *CLI) Run(args []string) int {
	flags := flag.NewFlagSet("notify_slack", flag.ContinueOnError)
	flags.SetOutput(c.errStream)

	var replaceExpr string

	flags.StringVar(&replaceExpr, "replace", "", `Replacement expression, e.g., "@search@replace@"`+"\n")

	err := flags.Parse(args[1:])
	if err != nil {
		return ExitCodeParseFlagError
	}

	if flags.NArg() == 0 {
		fmt.Fprintf(c.errStream, "Usage: purl --replace \"@search@replace@\" filename\n")
		return ExitCodeFail
	}
	filePath := flags.Arg(0)

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

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(c.errStream, "Failed to open file: %s\n", err)
		return ExitCodeFail
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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
		fmt.Printf("Error reading file: %s\n", err)
		return ExitCodeFail
	}

	return ExitCodeOK
}
