# Purl

Purl is a versatile text processing tool designed to easily and efficiently modify and replace text in files or from standard input. Inspired by the action of purling in knitting and the sound of a flowing stream, "Purl" symbolizes the concept of seamless repetition and smooth progress. Just as purling creates a fabric through consistent patterns and the stream's flow produces a calming rhythm, Purl facilitates effortless and repeated transformations of text. Aimed at providing the smoothness and efficiency of Perl one-liners, it is perfect for those looking for a tool to handle text processing tasks with ease and precision.

## Features

- **Auto Color Output by Default**: Purl automatically decides whether to colorize output based on the environment, enhancing readability. This auto-color feature aims to provide optimal visibility under various conditions without manual intervention.

- **Overwrite Option**: By specifying the `-overwrite` option, you can direct Purl to apply changes directly to the files. This functionality is not enabled by default to allow full control over when and how files are modified.

- **Flexible Input Options**: Purl accepts input either directly from specified files on the command line or through standard input, catering to a wide array of workflows and preferences.

## Options

- **`-overwrite`**: Use this option to enable Purl to overwrite the original files with the modified content. Without this option, Purl will display the results to standard output, leaving the original files unchanged.
- **`-replace`**: This option requires a replacement expression to specify the text you intend to change. Format your command as "@search@replace@", with "search" being the text to find and "replace" the text to insert.
- **`-color`** and **`-no-color`**: By default, Purl's output colorization is set to auto, determining the best mode based on your environment. Use `-no-color` if you prefer the output without colorization, regardless of the environment.
- **`-help`**: Display information about Purl and its various options.
- **`-version`**: Display version.

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

### Advanced Usage with Git Grep and Xargs

For users looking to apply replacements across multiple files in a Git repository:

```bash
git grep -l 'search_pattern' | xargs purl -overwrite -replace "@search_pattern@replace_text@"
```

This sequence finds all files containing 'search_pattern', then uses Purl to replace it with 'replace_text', directly modifying the files where the changes are applied.

Purl is crafted to offer simplicity for quick tasks as well as the capability to perform complex text processing, embodying the spirit of its name in every action it performs.
