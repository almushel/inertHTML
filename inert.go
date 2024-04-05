package main

import (
	"os"
	"strings"

	"github.com/almushel/inertHTML/nodes"
)

func GeneratePage(src, template, dest string) error {
	var srcFile, templateFile, destFile *os.File
	var err error

	srcFile, err = os.Open(src)
	defer srcFile.Close()
	if err != nil {
		return err
	}

	templateFile, err = os.Open(template)
	if err != nil {
		return err
	}

	destFile, err = CreateAll(dest)
	defer destFile.Close()
	if err != nil {
		return err
	}

	templateStr, err := ReadAllS(templateFile)
	if err != nil {
		return err
	}

	srcTxt, err := ReadAllS(srcFile)
	if err != nil {
		return err
	}

	blocks := nodes.ParseMDBlocks(srcTxt)
	blockNodes, err := nodes.BlocksToHTMLNodes(blocks)
	if err != nil {
		return err
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

	result := strings.Join(
		strings.Split(templateStr, "{{ Title }}"),
		pageTitle,
	)
	result = strings.Join(
		strings.Split(result, "{{ Content }}"),
		body,
	)

	return os.WriteFile(dest, []byte(result), 0666)
}

func GeneratePageRecursive(src, template, dest string) error {
	proc := func(path string) error {
		destPath := dest + path[len(src):]
		if strings.HasSuffix(path, ".md") {
			return GeneratePage(path, template, destPath[:len(destPath)-len(".md")]+".html")
		} else {
			return FileCopy(path, destPath)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return WalkDirRecursive(wd, src, proc)
}
