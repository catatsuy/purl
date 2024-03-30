package cli

func (c *CLI) ReplaceProcess(searchPattern string, replacement string) error {
	return c.replaceProcess(searchPattern, replacement)
}
