# inertHTML

A static site generator for converting markdown files to templated HTML.

## Usage

At a minimum, a positional source argument is required.
This can be either a single markdown file or a directory containing markdown files.
By default, the output will be written to equivalent html files in the same directory.

```sh
# Parses file.md and outputs file.html in the same directory
inertHTML file.md

# Parses *.md files in directory and outputs *.html files in place
inertHTML directory
```

### Output

If the `-o` flag is defined, the results will be output there.
Like the source, this can be a file or a directory.
If the given path does not exist, it is assumed to be a file only if the extension is `.html`.
Setting `-o` to a file when the source is a directory will return an error.

```sh
# Writes to file.html
inertHTML -o file.html file.md

# Writes to destDir/file.html
inertHTML -o destDir file.md

# Writes srcDir/*.md to destDir/*.html
inertHTML -o destDir srcDir

# Invalid combination. Returns an error.
inertHTML -o file.html srcDir
```

### Recursion

By default only the files at the top level of the source directory will be processed.
To recursively process all subdirectories, enable the `-r` flag.
The structure of the source directory will be reproduced in the output folder if `-o` is enabled.
This will be ignored if the source is a file.

```sh
# Parses *.md files in directory and all subdirectories and outputs *.html files in place
inertHTML -r directory
```

### Templates

A custom template file can be specified with the `-t` flag.
This file must be a valid html file with `<html>`, `<head>`, and `<body>` tags
as well as `{{ Title }}` and `{{ Content }}` template tags.

```sh
inertHTML -t template.html file.md
```

### Overwriting files

By default, inertHTML will quietly replace the contents of existing destination files.
This behavior can be changed with the following boolean flags:

* `-n`: No clobber. Quietly skips any existing files. 
* `-i`: Interactive mode. Asks for confirmation to overwrite each existing file.

## Markdown Features

inertHTML currently supports the majority of standard markdown syntax and some extensions,
prioritized based on what I needed for my own use.

### Limitations

Standard markdown syntax that is currently not supported (i.e. to-do):

- Nested lists
- Double trailing space line breaks
- Indent-based code blocks

### Extensions

- Fenced codeblocks
- Tables
- HTML in .md files
- Limited YAML frontmatter (detected and removed from output)

