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

	filePath    string
	replaceExpr string
	isOverwrite bool
	filters     rawStrings
	excludes    rawStrings
	help        bool
	color       bool
}

func NewCLI(outStream, errStream io.Writer, inputStream io.Reader) *CLI {
	return &CLI{outStream: outStream, errStream: errStream, inputStream: inputStream}
}

func (c *CLI) Run(args []string) int {
	flags, err := c.parseFlags(args)
	if err != nil {
		fmt.Fprintf(c.errStream, "Failed to parse flags: %s\n", err)
		return ExitCodeParseFlagError
	}

	if c.help {
		flags.Usage()
		return ExitCodeOK
	}

	err = c.validateInput(flags)
	if err != nil {
		fmt.Fprintf(c.errStream, "Failed to validate input: %s\n", err)
		return ExitCodeFail
	}

	if c.filePath != "" {
		file, err := os.Open(c.filePath)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to open file: %s\n", err)
			return ExitCodeFail
		}
		defer file.Close()

		c.inputStream = file
	}

	var tmpFile *os.File

	if c.isOverwrite {
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

	if len(c.replaceExpr) > 0 {
		delimiter := string(c.replaceExpr[0])
		parts := regexp.MustCompile(regexp.QuoteMeta(delimiter)).Split(c.replaceExpr[1:], -1)
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

	if len(c.filters) > 0 || len(c.excludes) > 0 {
		filters, err := compileRegexps(c.filters)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to compile regex patterns: %s\n", err)
			return ExitCodeFail
		}

		excludes, err := compileRegexps(c.excludes)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to compile regex patterns: %s\n", err)
			return ExitCodeFail
		}

		err = c.filterProcess(filters, excludes)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to process files: %s\n", err)
			return ExitCodeFail
		}
	}

	if c.isOverwrite {
		if err := os.Rename(tmpFile.Name(), c.filePath); err != nil {
			fmt.Fprintf(c.errStream, "Failed to overwrite the original file: %s\n", err)
			return ExitCodeFail
		}
	}

	return ExitCodeOK
}

func (c *CLI) parseFlags(args []string) (*flag.FlagSet, error) {
	flags := flag.NewFlagSet("purl", flag.ContinueOnError)
	flags.SetOutput(c.errStream)

	var noColor bool

	flags.BoolVar(&c.isOverwrite, "overwrite", false, "overwrite the file in place")
	flags.StringVar(&c.replaceExpr, "replace", "", `Replacement expression, e.g., "@search@replace@"`)
	flags.Var(&c.filters, "filter", `filter expression`)
	flags.Var(&c.excludes, "exclude", `exclude expression`)
	flags.BoolVar(&c.color, "color", false, `Colorize output`)
	flags.BoolVar(&noColor, "no-color", false, `Disable colorize output`)
	flags.BoolVar(&c.help, "help", false, `Show help`)

	flags.Usage = func() {
		fmt.Fprintln(c.errStream, "Usage: purl [options] [file]")
		flags.PrintDefaults()
	}

	err := flags.Parse(args[1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	c.color = !noColor && (c.color || term.IsTerminal(int(os.Stdout.Fd())))

	return flags, nil
}

func (c *CLI) validateInput(flags *flag.FlagSet) error {
	if flags.NArg() == 0 && term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("no input file specified")
	}

	if flags.NArg() == 0 && c.isOverwrite {
		return fmt.Errorf("cannot use -overwrite option with stdin")
	}

	if len(c.replaceExpr) != 0 && (len(c.filters) != 0 || len(c.excludes) != 0) {
		return fmt.Errorf("cannot use -replace and -filter options together")
	}

	if (len(c.filters) == 0 && len(c.excludes) == 0) && len(c.replaceExpr) < 3 {
		return fmt.Errorf("invalid replace expression format. Use \"@search@replace@\"")
	}

	if flags.NArg() == 0 {
		return nil
	}

	c.filePath = flags.Arg(0)
	if _, err := os.Stat(c.filePath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist")
	}

	return nil
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

func (c *CLI) filterProcess(filters []*regexp.Regexp, excludes []*regexp.Regexp) error {
	scanner := bufio.NewScanner(c.inputStream)

	for scanner.Scan() {
		line := scanner.Text()
		hit, hitRe := matchesFilters(line, filters)
		if len(filters) == 0 || hit {
			if excludeHit, _ := matchesFilters(line, excludes); !excludeHit {
				if hitRe != nil && c.color {
					line = colorText(line, hitRe)
				}
				fmt.Fprintln(c.outStream, line)
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

func matchesFilters(line string, regexps []*regexp.Regexp) (bool, *regexp.Regexp) {
	for _, re := range regexps {
		if re.MatchString(line) {
			return true, re
		}
	}
	return false, nil
}

func colorText(line string, re *regexp.Regexp) string {
	return re.ReplaceAllString(line, "\x1b[1m\x1b[91m$0\x1b[0m")
}
