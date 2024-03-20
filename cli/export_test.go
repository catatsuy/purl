package cli

import "io"

func (c *CLI) SetOutputStream(outStream io.Writer) {
	c.outStream = outStream
}

func (c *CLI) ProcessFiles(searchPattern string, replacement string) error {
	return c.processFiles(searchPattern, replacement)
}
