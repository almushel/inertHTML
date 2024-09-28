package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/almushel/inertHTML/parser"
)

type InertFlags struct {
	Interactive bool // Prompt stdin for file overwrites
	NoClobber   bool // Quietly skip any existing files
	Recursive   bool // Recersively process subdirectories
	Verbose     bool // Print steps to stdout
	PagesAsDirs bool // Output individual files to "filename/index.html", except for files named "index.md"
}

// Process markdown in src and output to dest using html template
// src: Path to markdown input file
// template: Path to html template file
// dest: Path to output html file
func GeneratePage(src, template, dest string) error {
	var err error
	var templateStr string

	if template != "" {
		templateStr, err = ReadFileS(template)
		if err != nil {
			return err
		}
	}

	srcTxt, err := ReadFileS(src)
	if err != nil {
		return err
	}

	// Detect and remove YAML frontmatter
	if strings.HasPrefix(srcTxt, "---\n") {
		_, body, found := strings.Cut(srcTxt[len("---\n"):], "\n---")
		if found {
			srcTxt = body
		}
	}
	// TODO: Process YAML frontmatter

	destFile, err := CreateAll(dest)
	defer destFile.Close()
	if err != nil {
		return err
	}

	result, err := parser.MDtoHTML(srcTxt)
	if err != nil {
		return err
	}

	return os.WriteFile(
		dest,
		[]byte(PopulateTemplate(result.Body, result.Title, templateStr)),
		0666,
	)
}

// Call generatePage with inert flag behaviors
func GeneratePageEx(src, template, dest string, flags InertFlags) error {
	if _, err := os.Stat(dest); !errors.Is(err, os.ErrNotExist) {
		if flags.NoClobber {
			return nil
		}

		if flags.Interactive {
			fmt.Printf("inertHTML: overwrite '%s'? ", dest)
			var input string
			_, err = fmt.Scanln(&input)
			if input != "y" && input != "yes" {
				return err
			}
		}
	}

	if flags.PagesAsDirs {
		info, err := os.Stat(src)
		if err != nil {
			return err
		}

		if info.Name() != "index.md" {
			dest = dest[:len(dest)-len(".html")] + "/index.html"
		}

	}

	if flags.Verbose {
		fmt.Printf("MD -> HTML: %s -> %s\n", src, dest)
	}
	return GeneratePage(src, template, dest)
}

// Process all markdown files in destination directory
// If recursive flag is set, continue recursively into subdirectories
// src:		Path to markdown input directory
// template:	Path to html template file
// dest:	Path to output directory
// flags:	Flags that modify generator behavior
func GenerateDirectory(src, template, dest string, flags InertFlags) error {
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	var srcPath, destPath string
	for _, file := range files {
		srcPath = filepath.Join(src, file.Name())
		destPath = filepath.Join(dest, file.Name())

		if flags.Recursive && file.IsDir() {
			if flags.Verbose {
				fmt.Printf("Processing directory: %s\n", srcPath)
			}
			err = GenerateDirectory(srcPath, template, destPath, flags)
		} else if filepath.Ext(srcPath) == ".md" {
			destFilePath := destPath[:len(destPath)-len("md")] + "html"
			err = GeneratePageEx(srcPath, template, destFilePath, flags)
		}
	}

	return err
}
