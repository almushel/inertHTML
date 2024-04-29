# inertHTML
A static site generator for converting markdown files to simple HTML.

## Usage

At a minimum, a positional source argument is required.
This can be either a single markdown file or a directory containing markdown files.
By default, the output will be written to quivalent html files in the same directory.

```sh
# Parses file.md and outputs file.html in the same directory
inertHTML file.md

# Parses *.md files in directory and outputs *.html files in place
inertHTML directory
```

### Recursion

By default only the files at the top level of the source directory will be processed.
To recursively process all subdirectories, enable the `-r` flag.
This will be ignored if the source is a file.

```sh
# Parses *.md files in directory and all subdirectories and outputs *.html files in place
inertHTML -r directory
```


### Output

If the `-o` flag is defined, the results will be output there.
Like the source, this can be a file or a directory.
However, setting `-o` to a directory for a file source will return an error.

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

### Templates

A custom template file can be specified with the `-t` flag.
This file must be a valid html file with `<html>` and `<body>` tags
as well as `{{ Title }}` and `{{ Content }}` template tags.

```sh
inertHTML -t template.html file.md
```

### Overwriting files

By default, inertHTML will quietly replace the contents of existing destination files.
This behavior can be changed with the following boolean flags:

* `-n`: No clobber. Quietly skips any existing files. 
* `-i`: Interactive mode. Asks for confirmation to overwrite each existing file.

