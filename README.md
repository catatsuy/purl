# Purl

Purl is a command-line utility designed for text file parsing and manipulation, offering a modern alternative to traditional sed and perl one-liners. It features intuitive options for filtering, transforming, and managing text data. Importantly, Purl accepts both file input and standard input, providing flexibility for various workflows. Moreover, it supports multiple instances of the -filter and -exclude options, allowing users to apply complex patterns of search refinements and exclusions in a single command.

Unlike sed, Purl works the same way on both Mac and Linux, without any compatibility issues. Simply download it to start using, offering a straightforward experience for text manipulation.

## Inspiration Behind "Purl"

The name "Purl" comes from a knitting technique and the sound of a flowing river. In knitting, "purling" is a way to create smooth fabric by repeating certain stitches. Similarly, the sound of a river is continuous and calming. "Purl" represents how this tool helps you transform text smoothly, just like the consistent pattern of knitting and the steady flow of a river. It is designed to use simple commands like Perl to make complex text tasks easy.

## Why purl?

Older tools like sed and grep have problems:

- **Complex Regex**: They use a regex that is not like Perl's, which is hard to use.
- **Not Compatible**: sed does not work the same on macOS and Linux, leading to problems.

I used Perl before because it is easier for regex and works well on different systems. But Perl has issues too:

- **Not Installed by Default**: New versions of macOS and Linux do not have Perl already installed, which makes it hard to use, especially in Docker.
- **Less Popular**: Fewer people use Perl now, so it seems less appealing.
- **Not Just for One-Liners**: Perl is not just for quick commands, which can make it confusing.

We need a new tool that:

- **Uses Easy Perl-Like Regex**: Easy to use for handling text using Go's regexp.
- **Works on All Systems**: No problems with different operating systems.
- **Light and Easy to Install**: Easy to set up anywhere.
- **Simple Commands**: Easy to understand, even for beginners.

So, I made **purl**. It brings together Perl's best parts but is simpler and ready for today's users.

## demo

This demo highlights how Purl utilizes both `-filter` and `-exclude` options to pinpoint essential lines within a log file, streamlining the search for relevant data.

https://github.com/catatsuy/purl/assets/1249910/72c01b33-082f-4b7f-84bc-4c59b0859df9

This demo shows how to use Purl to change `http://` URLs to `https://` in a source code file, using the `-replace` and `-overwrite` options to update the file directly.

https://github.com/catatsuy/purl/assets/1249910/be87d17a-44e7-4091-bc80-77921174eac2

This demo demonstrates using Purl to remove comments and empty lines from a configuration file, employing both `-filter` and `-exclude` options along with `-overwrite` to directly modify the file.

https://github.com/catatsuy/purl/assets/1249910/5cc479cc-ce1c-4901-864d-963bf659e125

## Features

**purl** is a tool that helps you easily handle data from different sources. Here are its main features:

- **Flexible Data Input and Output**: You can input data from typing directly or from files. Similarly, you can choose to output data directly to your screen or save it back to files.
- **Extract Specific Text**: The `-extract` option lets you extract and format specific parts of the text using regular expressions. For example, you can capture groups in the pattern and use them in a custom output format.
- **Simple Commands**: Use straightforward options like `-replace`, `-filter`, `-exclude`, and `-extract` to manage your data.
- **Edit Files Easily**: The `-overwrite` option lets you update files directly, making changes quick and simple.
- **Colorful Output**: When using the `-filter` option, the output on your screen can be colorful. You can control this with the `-color` or `-no-color` options.
- **Error on No Matches**: With the `-fail` option, Purl returns an error (status code 1) if no matches are found when using `-filter`, `-replace`, or `-extract`. If not used, Purl will not return an error even if no matches are found.

This tool is made to be user-friendly and effective for different data handling tasks.

## Input Processing Modes

Purl has two ways to read input:

1. **Multi-Line Mode (For Files)**
   When you use a file, purl reads the whole file at once. This lets you use regular expressions that match across multiple lines. However, it waits until the file is finished, so it is not good for live data (like piped input).

2. **Line Mode (For Standard Input)**
   When you do not give a file (using standard input), purl automatically uses **line mode**. In this mode, purl reads one line at a time and prints the output immediately. You can also force this mode by using the `-line` option.

   **Examples:**
   - To process data from a file or pipe and force line-by-line reading:
     ```bash
     cat yourfile.txt | purl -line -replace "@search@replace@"
     ```
   - If no file is given, purl uses line mode by default so it does not wait forever:
     ```bash
     purl -line -filter "error"
     ```

This way, purl outputs data quickly when reading live input, and still supports multi-line matching when reading files.

### Regular Expressions and Multi-Line Mode in Purl

Purl uses Go's `regexp` package with **Multi-Line Mode** (`(?m)`) always enabled. This makes line-based text processing intuitive and powerful.

#### What is Multi-Line Mode?

In Multi-Line Mode:
- `^` matches the start of a line.
- `$` matches the end of a line, just before the newline character (`\n`).

This mode is particularly useful when working with text files, as it aligns regular expressions with line-based operations.

Here’s a simplified and clear explanation for the README:

### Special Case: Empty Lines

When using `^$`, it matches completely empty lines with zero characters. However, lines that only contain a newline character (`\n`) are not considered "empty" in Go's `regexp`.

If you want to match and exclude lines that only contain a newline (`\n`), you should use `^\n$` instead.

#### Example

**Input:**
```
line1

line2
line3

```

**Command:**
```bash
purl -exclude "^\n$" file.txt
```

**Output:**
```
line1
line2
line3
```

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

### Why Use Multi-Line Mode?

Purl's default Multi-Line Mode (`(?m)`) simplifies operations like filtering, replacing, and excluding lines. You don't need to manually include `(?m)` in your patterns—it's always enabled for you.

#### Example Use Cases

1. **Filter Lines Starting with a Word**
   ```bash
   purl -filter "^START" file.txt
   ```

   Matches lines that start with `START`.

2. **Exclude Empty Lines**
   ```bash
   purl -exclude "^\n$" file.txt
   ```

   Removes lines that only contain a newline (`\n`).

3. **Replace Text on Specific Lines**
   ```bash
   purl -replace "@^line1@REPLACED@" file.txt
   ```

   Replaces lines starting with `line1` with `REPLACED`.

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

### Extracting Specific Parts of Text

The `-extract` option allows you to capture specific parts of the text using regular expressions and output them in a custom format.

```bash
purl -extract "@quick ([a-z]+) fox@animal: $1@" yourfile.txt
```

This command captures the word after "quick" and before "fox" and formats it as `animal: <captured_word>`. For example, given the input:

```
The quick brown fox jumps over the lazy dog
quick red fox
```

The output will be:

```
animal: brown
animal: red
```

You can also use multiple capture groups to format more complex outputs:

```bash
purl -extract "@(\\w+) (\\w+) (\\w+)@Match: $1, $2, $3@" yourfile.txt
```

With the input:

```
apple banana cherry
grape orange pineapple
```

The output will be:

```
Match: apple, banana, cherry
Match: grape, orange, pineapple
```

### Using the `-fail` Option

```bash
purl -fail -filter "error" yourlog.log
```

This command will search for lines containing "error" in `yourlog.log` and return an error (status code 1) if no matches are found. This behavior is useful for scripts where the absence of a match should trigger an error.

### Combining `-fail` with `-replace`

```bash
purl -fail -replace "@search@replace@" yourfile.txt
```

In this example, if no instances of "search" are found in `yourfile.txt`, Purl will return an error, indicating that no replacements were made.

These additions should provide a clear explanation of the new `-fail` option, helping users understand how to use it effectively in their workflows.

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

## Tips and Tricks

- **Using Special Characters**: Purl supports special characters like `\n` (newline), `\t` (tab), and `\r` (carriage return) in both patterns and replacements. For example:

```bash
purl -replace "@pattern@\nreplacement@" file.txt
```

This replaces occurrences of `pattern` followed by a newline and text in the file with replacement.

- **Handling Single Quotes**: When using single quotes or other special characters that conflict with shell syntax, you can combine different quoting styles. For example:

```bash
purl -replace '@pattern@'"'"'replacement'"'"'@' file.txt
```

Here, `'"'"'` is used to escape the single quote within the single-quoted string.

These tips can help you effectively use Purl in more complex scenarios while navigating shell-specific limitations.

## FAQ

### What can I do with regular expressions?

This tool uses Go's [`regexp` package](https://pkg.go.dev/regexp) directly. So, any pattern supported by Go's regular expressions can be used.

### Can I use characters other than '@' in the `-replace` option?

Yes, you can use different characters besides '@' for the `-replace` option. You just need to make sure the character you choose is not in your pattern or replacement.

If you want to use a different character, like '#', you can do it like this:

```bash
purl -replace "#pattern#replacement#" file.txt
```

Be sure your character is not part of your pattern or replacement text.

### Why Doesn't `-exclude '^$'` Work for Empty Lines?

In Go's `regexp`:

- `^$` matches lines that are completely empty (zero characters).
- Lines with just a newline (`\n`) contain one character and are not matched by `^$`.

To handle such lines, use `^\n$` to explicitly match them.
