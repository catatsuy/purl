package cli_test

import (
	"bytes"
	"testing"

	"github.com/catatsuy/purl/cli"
)

func TestProcessFiles_replace(t *testing.T) {
	outStream, errStream, inputStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	cl := cli.NewCLI(errStream, inputStream)
	cl.SetOutputStream(outStream)

	inputStream.WriteString("searchb searchc")

	err := cl.ProcessFiles("search", "replacement")

	if err != nil {
		t.Errorf("Error=%q", err)
	}

	expected := "replacementb replacementc\n"
	if outStream.String() != expected {
		t.Errorf("Output=%q, want %q; error: %q", outStream.String(), expected, errStream.String())
	}
}
