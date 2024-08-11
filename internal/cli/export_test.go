package cli

import (
	"io"
	"regexp"
)

func (c *CLI) ReplaceProcess(searchRe *regexp.Regexp, replacement []byte, inputStream io.Reader) (bool, error) {
	return c.replaceProcess(searchRe, replacement, inputStream)
}

func (c *CLI) FilterProcess(filters []*regexp.Regexp, notFilters []*regexp.Regexp, inputStream io.Reader) error {
	return c.filterProcess(filters, notFilters, inputStream)
}

func CompileRegexps(rawPatterns []string, ignoreCase bool) ([]*regexp.Regexp, error) {
	return compileRegexps(rawPatterns, ignoreCase)
}
