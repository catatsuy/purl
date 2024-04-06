package cli_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/catatsuy/purl/internal/cli"
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
			expected: "replacementa replacementb\nSearcha Searchb\n",
		},
		"provide multiple files for replace": {
			args:     []string{"purl", "-replace", "@search@replacement@", "testdata/test.txt", "testdata/testa.txt"},
			expected: "replacementa replacementb\nSearcha Searchb\nreplacementc replacementd\nnot not not\n",
		},
		"provide file for ignore case": {
			args:     []string{"purl", "-i", "-replace", "@search@replacement@", "testdata/test.txt"},
			expected: "replacementa replacementb\nreplacementa replacementb\n",
		},
		"provide multiple files for ignore case": {
			args:     []string{"purl", "-i", "-replace", "@search@replacement@", "testdata/test.txt", "testdata/testa.txt"},
			expected: "replacementa replacementb\nreplacementa replacementb\nreplacementc replacementd\nnot not not\n",
		},
		"provide stdin for ignore case": {
			args:     []string{"purl", "-i", "-replace", "@search@replacement@"},
			input:    "searcha Search\nsearchc Searchd\n",
			expected: "replacementa replacement\nreplacementc replacementd\n",
		},
		"provide multiple files for filter": {
			args:     []string{"purl", "-filter", "search", "testdata/test.txt", "testdata/testa.txt"},
			expected: "searcha searchb\nsearchc searchd\n",
		},
		"color text": {
			args:     []string{"purl", "-filter", "search", "-color"},
			input:    "searchb\nreplace\nsearchc",
			expected: "\x1b[1m\x1b[91msearch\x1b[0mb\n\x1b[1m\x1b[91msearch\x1b[0mc\n",
		},
		"no color text": {
			args:     []string{"purl", "-filter", "search", "-no-color"},
			input:    "searchb\nreplace\nsearchc",
			expected: "searchb\nsearchc\n",
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
			desc:         "version option",
			args:         []string{"purl", "-version"},
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
			desc:         "-exclude",
			args:         []string{"purl", "-exclude", "not filter"},
			expectedCode: 0,
		},
		{
			desc:         "multiple -exclude",
			args:         []string{"purl", "-exclude", "not filter", "-exclude", "not filter2"},
			expectedCode: 0,
		},
		{
			desc:         "provide -filter and -exclude",
			args:         []string{"purl", "-filter", "filter", "-exclude", "not filter2"},
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

func TestRun_successForOverwrite(t *testing.T) {
	testCases := []struct {
		desc         string
		args         []string
		filename     string
		expected     string
		expectedCode int
	}{
		{
			desc:         "-overwrite and -replace option",
			args:         []string{"purl", "-replace", "@search@replacement@", "-overwrite", "testdata/test_for_overwrite.txt"},
			filename:     "test_for_overwrite.txt",
			expected:     "replacemente replacementf\nnot not not\n",
			expectedCode: 0,
		},
		{
			desc:         "-overwrite and -filter option",
			args:         []string{"purl", "-filter", "search", "-overwrite", "testdata/test_for_overwrite.txt"},
			filename:     "test_for_overwrite.txt",
			expected:     "searche searchf\n",
			expectedCode: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testFilePath := "./testdata/" + tc.filename
			backupFilePath := testFilePath + ".bak"

			if err := copyFile(t, testFilePath, backupFilePath); err != nil {
				t.Fatalf("failed to backup the test file: %v", err)
			}

			defer func() {
				if err := os.Rename(backupFilePath, testFilePath); err != nil {
					t.Fatalf("failed to restore the test file: %v", err)
				}
			}()

			outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
			cl := cli.NewCLI(outStream, errStream, os.Stdin)

			if got := cl.Run(tc.args); got != tc.expectedCode {
				t.Fatalf("Expected exit code %d, but got %d; error: %q", tc.expectedCode, got, errStream.String())
			}

			b, err := os.ReadFile(testFilePath)
			if err != nil {
				t.Fatalf("failed to read the test file: %v", err)
			}

			if string(b) != tc.expected {
				t.Errorf("Output=%q, want %q; error: %q", string(b), tc.expected, errStream.String())
			}
		})
	}
}

func copyFile(t *testing.T, src, dst string) error {
	t.Helper()
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
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

func TestRun_failToProvideFiles(t *testing.T) {
	testCases := []struct {
		desc         string
		args         []string
		expectedCode int
	}{
		{
			desc:         "fail to provide -replace",
			args:         []string{"purl", "-replace", "search@replacement", "testdata/test.txt"},
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

func TestRun_failToProvideOverwriteAndStdin(t *testing.T) {
	testCases := []struct {
		desc         string
		args         []string
		expectedCode int
	}{
		{
			desc:         "fail to provide -replace",
			args:         []string{"purl", "-replace", "@search@replacement", "-overwrite"},
			expectedCode: 1,
		},
		{
			desc:         "fail to provide -replace",
			args:         []string{"purl", "-filter", "search", "-overwrite"},
			expectedCode: 1,
		},
		{
			desc:         "fail to provide -replace",
			args:         []string{"purl", "-exclude", "search", "-overwrite"},
			expectedCode: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			outStream, errStream, inputStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
			cl := cli.NewCLI(outStream, errStream, inputStream)

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

	err := cl.ReplaceProcess("search", "replacement", inputStream)

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

	err := cl.ReplaceProcess("search", "replacement", inputStream)

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
		name       string
		patterns   []string
		ignoreCase bool
		wantError  bool
	}{
		{
			name:      "ValidPatterns",
			patterns:  []string{"^test", "end$", "[0-9]+"},
			wantError: false,
		},
		{
			name:       "ValidPatterns for ignore case",
			patterns:   []string{"^test", "end$", "[0-9]+"},
			ignoreCase: true,
			wantError:  false,
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
			got, err := cli.CompileRegexps(tt.patterns, tt.ignoreCase)
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
		excludes   []string
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
			name:       "-exclude: SingleMatch",
			input:      "apple\nbanana\ncherry\n",
			excludes:   []string{"banana"},
			wantOutput: "apple\ncherry\n",
		},
		{
			name:       "-exclude: MultipleMatches",
			input:      "apple\nbanana\ncherry\n",
			excludes:   []string{"apple", "cherry"},
			wantOutput: "banana\n",
		},
		{
			name:       "-exclude: NoMatch",
			input:      "apple\nbanana\ncherry\n",
			excludes:   []string{"date"},
			wantOutput: "apple\nbanana\ncherry\n",
		},
		{
			name:       "-exclude: EmptyInput",
			input:      "",
			excludes:   []string{"apple"},
			wantOutput: "",
		},
		{
			name:       "provide filter and not filter",
			input:      "apple\nbanana\napple cherry\ncherry\n",
			filters:    []string{"apple"},
			excludes:   []string{"cherry"},
			wantOutput: "apple\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			outStream, errStream, inputStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
			cl := cli.NewCLI(outStream, errStream, inputStream)
			inputStream.WriteString(tt.input)

			filters, err := cli.CompileRegexps(tt.filters, false)
			if err != nil {
				t.Errorf("CompileRegexps() error = %v", err)
				return
			}

			excludes, err := cli.CompileRegexps(tt.excludes, false)
			if err != nil {
				t.Errorf("CompileRegexps() error = %v", err)
				return
			}

			err = cl.FilterProcess(filters, excludes, inputStream)
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
