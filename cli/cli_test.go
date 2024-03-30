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

func TestRun_successProcess(t *testing.T) {
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

func TestRun_success(t *testing.T) {
	testCases := []struct {
		desc         string
		args         []string
		expectedCode int
	}{
		{
			desc:         "help option",
			args:         []string{"purl", "-help"},
			expectedCode: 0,
		},
		{
			desc:         "filter",
			args:         []string{"purl", "-filter", "search"},
			expectedCode: 0,
		},
		{
			desc:         "multiple -filter",
			args:         []string{"purl", "-filter", "search", "-filter", "search2"},
			expectedCode: 0,
		},
		{
			desc:         "replace",
			args:         []string{"purl", "-replace", "@search@replace@"},
			expectedCode: 0,
		},
		{
			desc:         "-not-filter",
			args:         []string{"purl", "-not-filter", "not filter"},
			expectedCode: 0,
		},
		{
			desc:         "multiple -not-filter",
			args:         []string{"purl", "-not-filter", "not filter", "-not-filter", "not filter2"},
			expectedCode: 0,
		},
		{
			desc:         "provide -filter and -not-filter",
			args:         []string{"purl", "-filter", "filter", "-not-filter", "not filter2"},
			expectedCode: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
			cl := cli.NewCLI(outStream, errStream, os.Stdin)

			if got := cl.Run(tc.args); got != tc.expectedCode {
				t.Fatalf("Expected exit code %d, but got %d; error: %q", tc.expectedCode, got, errStream.String())
			}
		})
	}
}

func TestRun_failToProvideStdin(t *testing.T) {
	testCases := []struct {
		desc         string
		args         []string
		expectedCode int
	}{
		{
			desc:         "fail to provide -replace",
			args:         []string{"purl", "-replace", "search@replacement"},
			expectedCode: 1,
		},
		{
			desc:         "fail to provide -filter and -replace",
			args:         []string{"purl", "-filter", "aaa", "-replace", "@search@replacement@"},
			expectedCode: 1,
		},
		{
			desc:         "non-existent file with -replace",
			args:         []string{"purl", "-replace", "@search@replacement@", "testdata/noexist.txt"},
			expectedCode: 1,
		},
		{
			desc:         "non-existent file with -filter",
			args:         []string{"purl", "-filter", "aaaaa", "testdata/noexist.txt"},
			expectedCode: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
			cl := cli.NewCLI(outStream, errStream, os.Stdin)

			if got := cl.Run(tc.args); got != tc.expectedCode {
				t.Fatalf("Expected exit code %d, but got %d; error: %q", tc.expectedCode, got, errStream.String())
			}
		})
	}
}

func TestReplaceProcess_replace(t *testing.T) {
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

func TestReplaceProcess_noMatch(t *testing.T) {
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

func TestCompileRegexps(t *testing.T) {
	tests := []struct {
		name      string
		patterns  []string
		wantError bool
	}{
		{
			name:      "ValidPatterns",
			patterns:  []string{"^test", "end$", "[0-9]+"},
			wantError: false,
		},
		{
			name:      "InvalidPattern",
			patterns:  []string{"["},
			wantError: true,
		},
		{
			name:      "EmptyPattern",
			patterns:  []string{""},
			wantError: false,
		},
		{
			name:      "MixedValidAndInvalidPatterns",
			patterns:  []string{"^test", "["},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cli.CompileRegexps(tt.patterns)
			if tt.wantError {
				if err == nil {
					t.Errorf("%s: expected an error but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("%s: unexpected error: %v", tt.name, err)
				}
				if len(got) != len(tt.patterns) {
					t.Errorf("%s: expected %d compiled regexps, got %d", tt.name, len(tt.patterns), len(got))
				}
			}
		})
	}
}

func TestFilterProcess(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		filters    []string
		notFilters []string
		wantOutput string
	}{
		{
			name:       "SingleMatch",
			input:      "apple\nbanana\ncherry\n",
			filters:    []string{"banana"},
			wantOutput: "banana\n",
		},
		{
			name:       "MultipleMatches",
			input:      "apple\nbanana\ncherry\n",
			filters:    []string{"apple", "cherry"},
			wantOutput: "apple\ncherry\n",
		},
		{
			name:       "NoMatch",
			input:      "apple\nbanana\ncherry\n",
			filters:    []string{"date"},
			wantOutput: "",
		},
		{
			name:       "EmptyInput",
			input:      "",
			filters:    []string{"apple"},
			wantOutput: "",
		},
		{
			name:       "-not-filter: SingleMatch",
			input:      "apple\nbanana\ncherry\n",
			notFilters: []string{"banana"},
			wantOutput: "apple\ncherry\n",
		},
		{
			name:       "-not-filter: MultipleMatches",
			input:      "apple\nbanana\ncherry\n",
			notFilters: []string{"apple", "cherry"},
			wantOutput: "banana\n",
		},
		{
			name:       "-not-filter: NoMatch",
			input:      "apple\nbanana\ncherry\n",
			notFilters: []string{"date"},
			wantOutput: "apple\nbanana\ncherry\n",
		},
		{
			name:       "-not-filter: EmptyInput",
			input:      "",
			notFilters: []string{"apple"},
			wantOutput: "",
		},
		{
			name:       "provide filter and not filter",
			input:      "apple\nbanana\ncherry\n",
			filters:    []string{"apple"},
			notFilters: []string{"cherry"},
			wantOutput: "apple\nbanana\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			outStream, errStream, inputStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
			cl := cli.NewCLI(outStream, errStream, inputStream)
			inputStream.WriteString(tt.input)

			filters, err := cli.CompileRegexps(tt.filters)
			if err != nil {
				t.Errorf("CompileRegexps() error = %v", err)
				return
			}

			notFilters, err := cli.CompileRegexps(tt.notFilters)
			if err != nil {
				t.Errorf("CompileRegexps() error = %v", err)
				return
			}

			err = cl.FilterProcess(filters, notFilters)
			if err != nil {
				t.Errorf("filterProcess() error = %v", err)
				return
			}

			gotOutput := outStream.String()
			if gotOutput != tt.wantOutput {
				t.Errorf("filterProcess() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
