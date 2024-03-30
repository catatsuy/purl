package cli

import "regexp"

func (c *CLI) ReplaceProcess(searchPattern string, replacement string) error {
	return c.replaceProcess(searchPattern, replacement)
}

func (c *CLI) FilterProcess(regexps []*regexp.Regexp) error {
	return c.filterProcess(regexps)
}

func CompileRegexps(rawPatterns []string) ([]*regexp.Regexp, error) {
	return compileRegexps(rawPatterns)
}
