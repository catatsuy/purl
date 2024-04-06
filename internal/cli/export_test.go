package cli

import (
	"io"
	"regexp"
)

func (c *CLI) ReplaceProcess(searchPattern string, replacement string, inputStream io.Reader) error {
	return c.replaceProcess(searchPattern, replacement, inputStream)
}

func (c *CLI) FilterProcess(filters []*regexp.Regexp, notFilters []*regexp.Regexp, inputStream io.Reader) error {
	return c.filterProcess(filters, notFilters, inputStream)
}

func CompileRegexps(rawPatterns []string) ([]*regexp.Regexp, error) {
	return compileRegexps(rawPatterns)
}
