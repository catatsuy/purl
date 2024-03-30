package cli

import "regexp"

func (c *CLI) ReplaceProcess(searchPattern string, replacement string) error {
	return c.replaceProcess(searchPattern, replacement)
}

func (c *CLI) FilterProcess(filters []*regexp.Regexp, notFilters []*regexp.Regexp) error {
	return c.filterProcess(filters, notFilters)
}

func CompileRegexps(rawPatterns []string) ([]*regexp.Regexp, error) {
	return compileRegexps(rawPatterns)
}
