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

type rawStrings []string

func (i *rawStrings) String() string {
	return fmt.Sprint(*i)
}

func (i *rawStrings) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type CLI struct {
	outStream, errStream io.Writer
	inputStream          io.Reader
}

func NewCLI(outStream, errStream io.Writer, inputStream io.Reader) *CLI {
	return &CLI{outStream: outStream, errStream: errStream, inputStream: inputStream}
}

func (c *CLI) Run(args []string) int {
	flags := flag.NewFlagSet("purl", flag.ContinueOnError)
	flags.SetOutput(c.errStream)

	var replaceExpr string
	var isOverwrite bool
	var filters rawStrings

	flags.BoolVar(&isOverwrite, "overwrite", false, "overwrite the file in place")
	flags.StringVar(&replaceExpr, "replace", "", `Replacement expression, e.g., "@search@replace@"`)
	flags.Var(&filters, "filter", `Filter expression`)

	err := flags.Parse(args[1:])
	if err != nil {
		return ExitCodeParseFlagError
	}

	filePath := ""
	if flags.NArg() > 0 {
		filePath = flags.Arg(0)
	} else if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintln(c.errStream, "No input file specified")
		return ExitCodeFail
	}

	if isOverwrite && filePath == "" {
		fmt.Fprintln(c.errStream, "Cannot use -overwrite option with stdin")
		return ExitCodeFail
	}

	if len(replaceExpr) != 0 && len(filters) != 0 {
		fmt.Fprintln(c.errStream, "Cannot use -replace and -filter options together")
		return ExitCodeFail
	}

	if len(filters) == 0 && len(replaceExpr) < 3 {
		fmt.Fprintln(c.errStream, "Invalid replace expression format. Use \"@search@replace@\"")
		return ExitCodeFail
	}

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to open file: %s\n", err)
			return ExitCodeFail
		}
		defer file.Close()

		c.inputStream = file
	}

	var tmpFile *os.File

	if isOverwrite {
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
	}

	if len(replaceExpr) != 0 {
		delimiter := string(replaceExpr[0])
		parts := regexp.MustCompile(regexp.QuoteMeta(delimiter)).Split(replaceExpr[1:], -1)
		if len(parts) < 2 {
			fmt.Fprintln(c.errStream, "Invalid replace expression format. Use \"@search@replace@\"")
			return ExitCodeFail
		}
		searchPattern, replacement := parts[0], parts[1]

		if err := c.replaceProcess(searchPattern, replacement); err != nil {
			fmt.Fprintf(c.errStream, "Failed to process files: %s\n", err)
			return ExitCodeFail
		}
	}

	if len(filters) != 0 {
		regexps, err := compileRegexps(filters)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to compile regex patterns: %s\n", err)
			return ExitCodeFail
		}

		err = c.filterProcess(regexps)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to process files: %s\n", err)
			return ExitCodeFail
		}
	}

	if isOverwrite {
		if err := os.Rename(tmpFile.Name(), filePath); err != nil {
			fmt.Fprintf(c.errStream, "Failed to overwrite the original file: %s\n", err)
			return ExitCodeFail
		}
	}

	return ExitCodeOK
}

func (c *CLI) replaceProcess(searchPattern, replacement string) error {
	scanner := bufio.NewScanner(c.inputStream)

	re, err := regexp.Compile(searchPattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		modifiedLine := re.ReplaceAllString(line, replacement)
		fmt.Fprintln(c.outStream, modifiedLine)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func (c *CLI) filterProcess(regexps []*regexp.Regexp) error {
	scanner := bufio.NewScanner(c.inputStream)

	for scanner.Scan() {
		line := scanner.Text()

		for _, re := range regexps {
			if re.MatchString(line) {
				fmt.Fprintln(c.outStream, line)
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func compileRegexps(rawPatterns []string) ([]*regexp.Regexp, error) {
	regexps := make([]*regexp.Regexp, 0, len(rawPatterns))
	for _, pattern := range rawPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %w", err)
		}
		regexps = append(regexps, re)
	}
	return regexps, nil
}
