package main

import (
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

type InertFlags struct {
	NoClobber, Interactive bool
}

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

func GeneratePageEx(src, template, dest string, flags InertFlags) error {
	_, err := os.Stat(dest)
	if err == nil {
		if flags.NoClobber {
			return nil
		} else if flags.Interactive {
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

func GeneratePageRecursive(src, template, dest string) error {
	if s, _ := os.Stat(src); !s.IsDir() {
		return GeneratePage(src, template, dest)
	}

	proc := func(path string) error {
		destPath := dest + path[len(src):]
		if strings.HasSuffix(path, ".md") {
			destFilePath := destPath[:len(destPath)-len("md")] + "html"
			return GeneratePage(path, template, destFilePath)
		}

		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return WalkDirRecursive(wd, src, proc)
}

func GeneratePageRecursiveEx(src, template, dest string, flags InertFlags) error {
	if s, _ := os.Stat(src); !s.IsDir() {
		return GeneratePageEx(src, template, dest, flags)
	}

	proc := func(path string) error {
		destPath := dest + path[len(src):]
		if strings.HasSuffix(path, ".md") {
			destFilePath := destPath[:len(destPath)-len("md")] + "html"
			return GeneratePageEx(path, template, destFilePath, flags)
		}

		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return WalkDirRecursive(wd, src, proc)
}
