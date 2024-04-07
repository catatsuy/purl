# Purl

Purl is a command-line utility designed for text file parsing and manipulation, offering a modern alternative to traditional sed and perl one-liners. It features intuitive options for filtering, transforming, and managing text data. Importantly, Purl accepts both file input and standard input, providing flexibility for various workflows. Moreover, it supports multiple instances of the -filter and -exclude options, allowing users to apply complex patterns of search refinements and exclusions in a single command.

Unlike sed, Purl works the same way on both Mac and Linux, without any compatibility issues. Simply download it to start using, offering a straightforward experience for text manipulation.

## Inspiration Behind "Purl"

The name "Purl" is inspired by the dual notion of the knitting technique and the serene sound of a flowing stream. In knitting, purling refers to a method that creates a smooth, continuous fabric through repetitive patterns. Similarly, the tranquil sound of a stream embodies a steady, uninterrupted flow. "Purl" reflects the tool's ability to facilitate seamless and efficient text transformations, echoing the rhythmic repetition of knitting and the natural flow of water. Aimed at mirroring the simplicity and effectiveness of Perl one-liners, Purl is designed for those seeking a tool that combines ease of use with the precision needed for complex text processing tasks.

## Why Choose Purl?

Purl is a handy tool that can replace `sed` and `grep` for many tasks. It uses Go's regexp for regular expressions, which offers similar capabilities to PCRE but is not fully compatible. This means you get a powerful tool for text processing that works well on both Mac and Linux without worrying about differences between these systems. Moreover, Purl is easy to start using: just download and run. It's designed to be simple and efficient for anyone who needs to manipulate text.

## Features

- **Accepts Standard Input**: In addition to processing files, Purl can take input piped from other commands, expanding its usability.
- **Flexible Exclusions**: Skip lines matching specified regular expressions, focusing on relevant data.
- **Search Refinement**: Apply additional filters to your search queries for more precise results.
- **Case Insensitivity**: Perform case-insensitive searches with a simple option, broadening your search capabilities.
- **In-Place Editing**: Directly modify the original files with your changes, simplifying the workflow.
- **Colored Output**: Enhance readability with optional colored output, automatically adjusted based on your terminal's capabilities.
- **Custom Replacements**: Define custom replacement patterns for comprehensive text manipulation.

## Installation

It is recommended that you use the binaries available on [GitHub Releases](https://github.com/catatsuy/purl/releases). It is advisable to download and use the latest version.

If you have a Go language development environment set up, you can also compile and install the 'purl' tools on your own.

```bash
go install github.com/catatsuy/purl@latest
```

To build and modify the 'purl' tools for development purposes, you can use the `make` command.

```bash
make
```

If you use the `make` command to build and install the 'purl' tool, the output of the `purl -version` command will be the git commit ID of the current version.

### Using Purl with GitHub Actions

If you want to use Purl in your GitHub Actions workflows, include the following steps in your `.github/workflows` YAML file:

```yaml
- name: Download purl
  run: |
    curl -sL https://github.com/catatsuy/purl/releases/latest/download/purl-linux-amd64.tar.gz | tar xz -C /tmp

- name: Move purl to /usr/local/bin
  run: |
    sudo mv /tmp/purl /usr/local/bin/
```

These steps ensure that Purl is downloaded and moved to `/usr/local/bin`, making it available for use in subsequent steps of your workflow.

## Usage Examples

### Preview Changes Before Applying

```bash
purl -replace "@search@replace@" yourfile.txt
```

This command searches for "search" in `yourfile.txt`, shows how it would be replaced with "replace", but does not modify the file itself.

### Directly Modify Files

```bash
purl -overwrite -replace "@search@replace@" yourfile.txt
```

Using the `-overwrite` option, Purl will replace "search" with "replace" in `yourfile.txt` and save the changes to the file.

### Using Standard Input

Purl can also process input piped from other commands, offering flexibility in how it's used:

```bash
cat yourfile.txt | purl -replace "@search@replace@"
```

This feeds the content of `yourfile.txt` into Purl, which processes and displays the modified text according to the specified replacement pattern.

### Using multiple files

Purl supports processing multiple files in a single command, allowing you to apply operations across several documents simultaneously. Simply list the files at the end of your command. For example:

```bash
purl -replace "@search@replacement@" file1.txt file2.txt file3.txt
```

This command will apply the replace operation to 'search' with 'replacement' in `file1.txt`, `file2.txt`, and `file3.txt`.

### Usage with `-filter`

```bash
purl -filter "error" yourlog.log
```

This command filters the lines containing "error" in `yourlog.log`, displaying them with colored output for better visibility.

### Filtering Input with Multiple Criteria

To filter lines that meet multiple criteria, you can use the `-filter` option multiple times. This works both when reading from a file and processing standard input.

```bash
purl -filter "error" -filter "warning" yourlog.log
```

Or for standard input:

```bash
cat yourlog.log | purl -filter "error" -filter "warning"
```

This will display lines that contain either "error" or "warning" from `yourlog.log`.

### Excluding Lines with Multiple Patterns

Similarly, you can exclude lines that match multiple patterns by specifying `-exclude` more than once:

```bash
purl -exclude "debug" -exclude "info" yourlog.log
```

Or for piped input:

```bash
cat yourlog.log | purl -exclude "debug" -exclude "info"
```

This approach excludes lines that contain "debug" or "info" from the output.

Purl allows combining `-filter` and `-exclude` for precise text control.

### Using the -i Option for Case-Insensitive Searches

When the `-i` option is used with Purl, it allows case-insensitive matching for filters and exclusions. For instance:

```bash
purl -i -filter "error" yourfile.txt
```

This command will match lines containing 'error' in any case variation, such as 'Error', 'ERROR', or 'error', in `yourfile.txt`.

The `-i` option in Purl enables case-insensitive operations not only for `-filter` but also for `-exclude` and `-replace`. This enhances flexibility in handling text variations. For example:

```bash
purl -i -exclude "debug" yourfile.txt
```

This will exclude lines with 'debug', 'Debug', 'DEBUG', etc., from `yourfile.txt`.

Similarly, when using `-replace`:

```bash
purl -i -replace "@search@replacement@" yourfile.txt
```

This applies the replacement operation regardless of case differences between 'search' and its occurrences in `yourfile.txt`, ensuring 'Search', 'SEARCH', etc., are also matched and replaced.

### Integrating with Git, Grep, and Xargs

For users looking to apply replacements across multiple files in a Git repository:

```bash
git grep -l 'search_pattern' | xargs purl -overwrite -replace "@search_pattern@replace_text@"
```

This sequence finds all files containing 'search_pattern', then uses Purl to replace it with 'replace_text', directly modifying the files where the changes are applied.

Purl is crafted to offer simplicity for quick tasks as well as the capability to perform complex text processing, embodying the spirit of its name in every action it performs.
