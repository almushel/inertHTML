package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/almushel/inertHTML/parser"
)

const defaultTemplate = `<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title> {{ Title }} </title>
    <link href="/index.css" rel="stylesheet">
</head>

<body>
    <article>
        {{ Content }}
    </article>
</body>

</html>`

func ValidateTemplateFile(file string) error {
	info, err := os.Stat(file)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open template file: %s", err.Error()))
	}

	if info.IsDir() {
		return errors.New("Template file is a directory")
	}

	template, err := ReadFileS(file)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read template file: %s", err.Error()))
	}

	substrs := []string{
		"<html>", "</html>",
		"<body>", "</body>",
		"{{ Title }}", "{{ Content }}",
	}

	for _, s := range substrs {
		if !strings.Contains(template, s) {
			return errors.New(fmt.Sprintf("Invalid template: %s tag not found.", s))
		}
	}

	return nil
}

// Process markdown in src and output to dest using html template
func GeneratePage(src, template, dest string) error {
	var err error
	var templateStr string

	if template == "" {
		templateStr = defaultTemplate
	} else {
		templateStr, err = ReadFileS(template)
		if err != nil {
			return err
		}
	}

	srcTxt, err := ReadFileS(src)
	if err != nil {
		return err
	}

	destFile, err := CreateAll(dest)
	defer destFile.Close()
	if err != nil {
		return err
	}

	result, err := parser.MDtoHTML(srcTxt, templateStr)
	if err != nil {
		return err
	}

	fmt.Printf("MD -> HTML: %s -> %s\n", src, dest)
	return os.WriteFile(dest, []byte(result), 0666)
}

// GeneratePage with inert flag behaviors
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

	return GeneratePage(src, template, dest)
}

// Process all md files at dest
// If recursive flag is set, continue recursively into subdirectories
func GenerateDirectory(src, template, dest string, flags InertFlags) error {
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	var srcPath, destPath string
	for _, file := range files {
		srcPath = src + "/" + file.Name()
		destPath = dest + "/" + file.Name()

		if src == dest {
			continue
		}

		if flags.Recursive && file.IsDir() {
			err = GenerateDirectory(srcPath, template, dest, flags)
		} else if strings.HasSuffix(srcPath, ".md") {
			destFilePath := destPath[:len(destPath)-len("md")] + "html"
			err = GeneratePageEx(srcPath, template, destFilePath, flags)
		}
	}

	return err
}
