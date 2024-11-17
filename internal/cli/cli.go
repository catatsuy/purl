package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
)

const (
	ExitCodeOK             = 0
	ExitCodeParseFlagError = 2
	ExitCodeFail           = 2
	ExitCodeNoMatch        = 1
)

var (
	Version string
)

func version() string {
	if Version != "" {
		return Version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}
	return info.Main.Version
}

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

	isStdinTerminal  bool
	isStdoutTerminal bool

	filePaths   []string
	replaceExpr string
	isOverwrite bool
	filters     rawStrings
	excludes    rawStrings
	extractExpr string
	help        bool
	isColor     bool
	ignoreCase  bool
	failMode    bool
	version     bool

	appVersion string
}

func NewCLI(outStream, errStream io.Writer, inputStream io.Reader, isStdinTerminal, isStdoutTerminal bool) *CLI {
	return &CLI{appVersion: version(), outStream: outStream, errStream: errStream, inputStream: inputStream, isStdinTerminal: isStdinTerminal, isStdoutTerminal: isStdoutTerminal}
}

func unescapeString(input string) string {
	replacer := strings.NewReplacer(
		`\\`, `\`,
		`\n`, "\n",
		`\t`, "\t",
		`\r`, "\r",
	)
	return replacer.Replace(input)
}

func (c *CLI) Run(args []string) int {
	flags, err := c.parseFlags(args)
	if err != nil {
		fmt.Fprintf(c.errStream, "Failed to parse flags: %s\n", err)
		return ExitCodeParseFlagError
	}

	if c.version {
		fmt.Fprintf(c.errStream, "purl version %s; %s\n", c.appVersion, runtime.Version())
		return ExitCodeOK
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

	var searchRe *regexp.Regexp
	var replacement []byte

	if len(c.replaceExpr) > 0 {
		delimiter := string(c.replaceExpr[0])
		parts := regexp.MustCompile(regexp.QuoteMeta(delimiter)).Split(c.replaceExpr[1:], -1)
		if len(parts) < 2 {
			fmt.Fprintln(c.errStream, "Invalid replace expression format. Use \"@search@replace@\"")
			return ExitCodeFail
		}
		searchPattern := parts[0]
		replacementStr := unescapeString(parts[1])
		replacement = []byte(replacementStr)

		if c.ignoreCase {
			searchPattern = "(?i)" + searchPattern
		}

		searchRe, err = regexp.Compile(searchPattern)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to compile regex pattern: %s\n", err)
			return ExitCodeFail
		}
	}

	var filterRes, excludeRes []*regexp.Regexp

	if len(c.filters) > 0 || len(c.excludes) > 0 {
		filterRes, err = compileRegexps(c.filters, c.ignoreCase)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to compile regex patterns: %s\n", err)
			return ExitCodeFail
		}

		excludeRes, err = compileRegexps(c.excludes, c.ignoreCase)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to compile regex patterns: %s\n", err)
			return ExitCodeFail
		}
	}

	if len(c.filePaths) != 0 {
		for _, filePath := range c.filePaths {

			file, err := os.Open(filePath)
			if err != nil {
				fmt.Fprintf(c.errStream, "Failed to open file: %s\n", err)
				return ExitCodeFail
			}
			defer file.Close()

			var tmpFile *os.File

			if c.isOverwrite {
				tmpFile, err = os.CreateTemp("", "purl")
				if err != nil {
					fmt.Fprintf(c.errStream, "Failed to create temp file: %s\n", err)
					return ExitCodeFail
				}
				defer os.Remove(tmpFile.Name())
				defer tmpFile.Close()

				c.outStream, err = os.Create(tmpFile.Name())
				if err != nil {
					fmt.Fprintf(c.errStream, "Failed to open file for writing: %s\n", err)
					return ExitCodeFail
				}
			}

			if len(c.replaceExpr) > 0 {
				matched, err := c.replaceProcess(searchRe, replacement, file)
				if err != nil {
					fmt.Fprintf(c.errStream, "Failed to process files: %s\n", err)
					return ExitCodeFail
				}

				if c.failMode && !matched {
					fmt.Fprintf(c.errStream, "No matches found in file: %s\n", filePath)
					return ExitCodeNoMatch
				}
			}

			if len(c.filters) > 0 || len(c.excludes) > 0 {
				matched, err := c.filterProcess(filterRes, excludeRes, file)
				if err != nil {
					fmt.Fprintf(c.errStream, "Failed to process files: %s\n", err)
					return ExitCodeFail
				}

				if c.failMode && !matched {
					fmt.Fprintf(c.errStream, "No matches found in file: %s\n", filePath)
					return ExitCodeNoMatch
				}
			}

			if c.isOverwrite {
				if err := os.Rename(tmpFile.Name(), filePath); err != nil {
					fmt.Fprintf(c.errStream, "Failed to overwrite the original file: %s\n", err)
					return ExitCodeFail
				}
			}
		}
	} else {
		if len(c.replaceExpr) > 0 {
			matched, err := c.replaceProcess(searchRe, replacement, c.inputStream)
			if err != nil {
				fmt.Fprintf(c.errStream, "Failed to process files: %s\n", err)
				return ExitCodeFail
			}

			if c.failMode && !matched {
				fmt.Fprintln(c.errStream, "No matches found in input")
				return ExitCodeNoMatch
			}
		}

		if len(c.filters) > 0 || len(c.excludes) > 0 {
			matched, err := c.filterProcess(filterRes, excludeRes, c.inputStream)
			if err != nil {
				fmt.Fprintf(c.errStream, "Failed to process files: %s\n", err)
				return ExitCodeFail
			}

			if c.failMode && !matched {
				fmt.Fprintln(c.errStream, "No matches found in input")
				return ExitCodeNoMatch
			}
		}
	}

	if len(c.extractExpr) > 0 {
		// Split the extract expression into pattern and replacement
		delimiter := string(c.extractExpr[0])
		parts := regexp.MustCompile(regexp.QuoteMeta(delimiter)).Split(c.extractExpr[1:], -1)
		if len(parts) < 2 {
			fmt.Fprintln(c.errStream, "Invalid extract expression format. Use \"@pattern@replacement@\"")
			return ExitCodeFail
		}
		searchPattern := parts[0]
		replacementStr := unescapeString(parts[1])
		replacement := []byte(replacementStr)

		// Add case-insensitive flag if necessary
		if c.ignoreCase {
			searchPattern = "(?i)" + searchPattern
		}

		// Compile the regex pattern
		searchRe, err := regexp.Compile(searchPattern)
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed to compile extract regex pattern: %s\n", err)
			return ExitCodeFail
		}

		// Process files if provided
		if len(c.filePaths) > 0 {
			for _, filePath := range c.filePaths {
				file, err := os.Open(filePath)
				if err != nil {
					fmt.Fprintf(c.errStream, "Failed to open file: %s\n", err)
					return ExitCodeFail
				}
				defer file.Close()

				// Process the file content with the extract logic
				matched, err := c.extractProcess(searchRe, replacement, file)
				if err != nil {
					fmt.Fprintf(c.errStream, "Failed to process file: %s\n", err)
					return ExitCodeFail
				}

				// Handle fail mode if no matches are found
				if c.failMode && !matched {
					fmt.Fprintf(c.errStream, "No matches found in file: %s\n", filePath)
					return ExitCodeNoMatch
				}
			}
		} else {
			// Process standard input if no files are provided
			matched, err := c.extractProcess(searchRe, replacement, c.inputStream)
			if err != nil {
				fmt.Fprintf(c.errStream, "Failed to process input: %s\n", err)
				return ExitCodeFail
			}

			// Handle fail mode if no matches are found
			if c.failMode && !matched {
				fmt.Fprintln(c.errStream, "No matches found in input")
				return ExitCodeNoMatch
			}
		}

		return ExitCodeOK
	}

	return ExitCodeOK
}

func (c *CLI) parseFlags(args []string) (*flag.FlagSet, error) {
	flags := flag.NewFlagSet("purl", flag.ContinueOnError)
	flags.SetOutput(c.errStream)

	var color, noColor bool

	flags.BoolVar(&c.isOverwrite, "overwrite", false, "Replace original file with results.")
	flags.StringVar(&c.replaceExpr, "replace", "", "Format: '@match@replacement@'.")
	flags.StringVar(&c.extractExpr, "extract", "", "Extract and print text matching the regex pattern.")
	flags.Var(&c.filters, "filter", "Apply search refinement.")
	flags.Var(&c.excludes, "exclude", "Exclude lines matching regex.")
	flags.BoolVar(&color, "color", false, "Colored output. Default auto.")
	flags.BoolVar(&noColor, "no-color", false, "Disable colored output.")
	flags.BoolVar(&c.ignoreCase, "i", false, `Ignore case (prefixes '(?i)' to all regular expressions)`)
	flags.BoolVar(&c.failMode, "fail", false, "Exit with a non-zero status if no matches are found")
	flags.BoolVar(&c.help, "help", false, `Show help`)
	flags.BoolVar(&c.version, "version", false, "Print version and quit")

	flags.Usage = func() {
		fmt.Fprintf(c.errStream, "purl version %s; %s\nUsage: purl [options] [file]\n", c.appVersion, runtime.Version())
		flags.PrintDefaults()
	}

	err := flags.Parse(args[1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	c.isColor = !noColor && (color || c.isStdoutTerminal)

	return flags, nil
}

func (c *CLI) validateInput(flags *flag.FlagSet) error {
	if flags.NArg() == 0 && c.isStdinTerminal {
		return fmt.Errorf("no input file specified")
	}

	if flags.NArg() == 0 && c.isOverwrite {
		return fmt.Errorf("cannot use -overwrite option with stdin")
	}

	err := c.validateMutuallyExclusiveOptions()
	if err != nil {
		return err
	}

	if err := c.validateExpressionFormats(); err != nil {
		return err
	}

	if flags.NArg() > 0 {
		c.filePaths = flags.Args()
		for _, filePath := range c.filePaths {
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("input file does not exist: %s", filePath)
			}
		}
	}

	return nil
}

// validateMutuallyExclusiveOptions checks that incompatible options are not used together
func (c *CLI) validateMutuallyExclusiveOptions() error {
	if len(c.extractExpr) > 0 {
		if len(c.replaceExpr) > 0 || len(c.filters) > 0 || len(c.excludes) > 0 {
			return fmt.Errorf("-extract cannot be used with -replace, -filter, or -exclude options")
		}
	}

	if len(c.replaceExpr) > 0 && (len(c.filters) > 0 || len(c.excludes) > 0 || len(c.extractExpr) > 0) {
		return fmt.Errorf("-replace cannot be used with -filter, -exclude, or -extract options")
	}

	return nil
}

// validateExpressionFormats checks the format of expressions
func (c *CLI) validateExpressionFormats() error {
	// Validate -replace expression format
	if len(c.replaceExpr) > 0 && len(c.replaceExpr) < 3 {
		return fmt.Errorf("invalid replace expression format. Use \"@search@replace@\"")
	}

	// Validate -extract expression format
	if len(c.extractExpr) > 0 && len(c.extractExpr) < 3 {
		return fmt.Errorf("invalid extract expression format. Use \"@search@replace@\"")
	}

	return nil
}

// replaceProcess reads data from inputStream, performs a regex replacement,
// and writes the modified data to outputStream.
// If input is from a pipe, it processes input line by line without changing newline characters.
// If input is from a file, it reads and processes the entire file at once.
func (c *CLI) replaceProcess(searchRe *regexp.Regexp, replacement []byte, inputStream io.Reader) (bool, error) {
	matched := false
	if c.isStdinTerminal {
		// Read all data from the file input
		b, err := io.ReadAll(inputStream)
		if err != nil {
			return false, fmt.Errorf("error reading file: %w", err)
		}

		modified := searchRe.ReplaceAllFunc(b, func(match []byte) []byte {
			matched = true
			return replacement
		})
		c.outStream.Write(modified)
	} else {
		// Read input line by line when input is from a pipe without changing newline characters
		reader := bufio.NewReader(inputStream)
		for {
			line, err := reader.ReadBytes('\n')

			if err != nil && !errors.Is(err, io.EOF) {
				return false, fmt.Errorf("error reading input: %w", err)
			}

			if errors.Is(err, io.EOF) && len(line) == 0 {
				break
			}

			// Replace text in each line using the regex
			modifiedLine := searchRe.ReplaceAllFunc(line, func(match []byte) []byte {
				matched = true
				return replacement
			})

			// Write the changed line to the output
			if _, err := c.outStream.Write(modifiedLine); err != nil {
				return false, fmt.Errorf("error writing to output: %w", err)
			}
		}
	}

	return matched, nil
}

func (c *CLI) filterProcess(filters []*regexp.Regexp, excludes []*regexp.Regexp, inputStream io.Reader) (bool, error) {
	matched := false
	// Read input line by line when input is from a pipe without changing newline characters
	reader := bufio.NewReader(inputStream)
	for {
		line, err := reader.ReadBytes('\n')

		if err != nil && !errors.Is(err, io.EOF) {
			return false, fmt.Errorf("error reading input: %w", err)
		}

		if errors.Is(err, io.EOF) && len(line) == 0 {
			break
		}

		hit, hitRes := matchesFilters(line, filters)
		if len(filters) == 0 || hit {
			matched = hit
			if excludeHit, _ := matchesFilters(line, excludes); !excludeHit {
				if len(hitRes) > 0 && c.isColor {
					line = colorText(line, hitRes)
				}

				if _, err := c.outStream.Write(line); err != nil {
					return false, fmt.Errorf("error writing to output: %w", err)
				}
			}
		}
	}

	return matched, nil
}

func (c *CLI) extractProcess(searchRe *regexp.Regexp, replacement []byte, inputStream io.Reader) (bool, error) {
	matched := false
	replacementStr := string(replacement)

	b, err := io.ReadAll(inputStream)
	if err != nil {
		return false, fmt.Errorf("error reading file: %w", err)
	}

	matches := searchRe.FindAllSubmatch(b, -1)
	for _, match := range matches {
		matched = true

		// Construct replacements for placeholders dynamically
		replacements := make([]string, 0, 2*len(match))
		for i := len(match) - 1; i >= 0; i-- { // Start from the largest index
			replacements = append(replacements, fmt.Sprintf("$%d", i), string(match[i]))
		}

		// Apply the replacements
		result := strings.NewReplacer(replacements...).Replace(replacementStr)

		// Write the result with error checking
		if _, err := fmt.Fprintln(c.outStream, result); err != nil {
			return false, fmt.Errorf("error writing to output: %w", err)
		}
	}

	return matched, nil
}

func compileRegexps(rawPatterns []string, ignoreCase bool) ([]*regexp.Regexp, error) {
	regexps := make([]*regexp.Regexp, 0, len(rawPatterns))
	for _, pattern := range rawPatterns {
		if ignoreCase {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %w", err)
		}
		regexps = append(regexps, re)
	}
	return regexps, nil
}

func matchesFilters(line []byte, regexps []*regexp.Regexp) (bool, []*regexp.Regexp) {
	var matchedRegexps []*regexp.Regexp
	for _, re := range regexps {
		if re.Match(line) {
			matchedRegexps = append(matchedRegexps, re)
		}
	}
	return len(matchedRegexps) > 0, matchedRegexps
}

func colorText(line []byte, res []*regexp.Regexp) []byte {
	for _, re := range res {
		line = re.ReplaceAll(line, []byte("\x1b[1m\x1b[91m$0\x1b[0m"))
	}
	return line
}
