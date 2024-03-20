package cli

func (c *CLI) ProcessFiles(searchPattern string, replacement string) error {
	return c.processFiles(searchPattern, replacement)
}
