package cli_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/catatsuy/purl/cli"
)

func TestNewCLI(t *testing.T) {
	outStream, errStream, inputStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	cl := cli.NewCLI(outStream, errStream, inputStream)

	if cl == nil {
		t.Error("NewCLI should not return nil")
	}
}

func TestRun_success(t *testing.T) {
	tests := map[string]struct {
		args     []string
		input    string
		expected string
		result   int
	}{
		"normal": {
			args:     []string{"purl", "-replace", "@search@replacement@"},
			input:    "searchb searchc",
			expected: "replacementb replacementc\n",
		},
		"no match": {
			args:     []string{"purl", "-replace", "@search@replacement@"},
			input:    "no match",
			expected: "no match\n",
		},
		"provide file": {
			args:     []string{"purl", "-replace", "@search@replacement@", "testdata/test.txt"},
			expected: "replacementa replacementb\n",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			outStream, errStream, inputStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
			cl := cli.NewCLI(outStream, errStream, inputStream)
			inputStream.WriteString(test.input)

			expectedCode := 0
			if got, expected := cl.Run(test.args), expectedCode; got != expected {
				t.Fatalf("Expected exit code %d, but got %d; error: %q", expected, got, errStream.String())
			}

			if outStream.String() != test.expected {
				t.Errorf("Output=%q, want %q; error: %q", outStream.String(), test.expected, errStream.String())
			}
		})
	}
}

func TestRun_failToProvideStdin(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cl := cli.NewCLI(outStream, errStream, os.Stdin)

	expectedCode := 1
	if got, expected := cl.Run([]string{"purl", "-replace", "@search@replacement@"}), expectedCode; got == expected {
		t.Fatalf("Expected exit code %d, but got %d; error: %q", expected, got, errStream.String())
	}
}

func TestProcessFiles_replace(t *testing.T) {
	outStream, errStream, inputStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	cl := cli.NewCLI(outStream, errStream, inputStream)

	inputStream.WriteString("searchb searchc")

	err := cl.ReplaceProcess("search", "replacement")

	if err != nil {
		t.Errorf("Error=%q", err)
	}

	expected := "replacementb replacementc\n"
	if outStream.String() != expected {
		t.Errorf("Output=%q, want %q; error: %q", outStream.String(), expected, errStream.String())
	}
}

func TestProcessFiles_noMatch(t *testing.T) {
	outStream, errStream, inputStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	cl := cli.NewCLI(outStream, errStream, inputStream)

	inputStream.WriteString("no match")

	err := cl.ReplaceProcess("search", "replacement")

	if err != nil {
		t.Errorf("Error=%q", err)
	}

	expected := "no match\n"
	if outStream.String() != expected {
		t.Errorf("Output=%q, want %q; error: %q", outStream.String(), expected, errStream.String())
	}
}
