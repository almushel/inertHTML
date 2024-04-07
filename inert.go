package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/almushel/inertHTML/nodes"
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

func MDtoHTML(src, template string) (string, error) {
	var result string

	blocks := nodes.ParseMDBlocks(src)
	blockNodes, err := nodes.BlocksToHTMLNodes(blocks)
	if err != nil {
		return result, err
	}

	var pageTitle string

	for i := range blockNodes {
		blockNodes[i].ProcessInnerText()
	}

	var body string
	for _, node := range blockNodes {
		if pageTitle == "" && node.Tag == "h1" {
			pageTitle = node.Value
		}
		body += node.ToHTML()
	}

	result = strings.Join(
		strings.Split(template, "{{ Title }}"),
		pageTitle,
	)
	result = strings.Join(
		strings.Split(result, "{{ Content }}"),
		body,
	)

	return result, nil
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

	result, err := MDtoHTML(srcTxt, templateStr)
	if err != nil {
		return err
	}

	fmt.Printf("MD -> HTML: %s -> %s\n", src, dest)
	return os.WriteFile(dest, []byte(result), 0666)
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
