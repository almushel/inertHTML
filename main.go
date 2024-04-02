package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// os.Create for complete path to file
func CreateAll(name string) (*os.File, error) {
	var i int = strings.LastIndex(name, "/")
	if i >= 0 {
		err := os.MkdirAll(name[:i], 0700)
		if err != nil {
			return nil, err
		}
	}

	return os.Create(name)
}

func ReadAll(src *os.File) ([]byte, error) {
	var result []byte
	var err error

	var buf []byte = make([]byte, 2048)
	var bytesRead int
	for bytesRead, err = src.Read(buf); bytesRead != 0; bytesRead, err = src.Read(buf) {
		result = append(result, buf[:bytesRead]...)
	}

	if err == io.EOF {
		err = nil
	}

	return result, err
}

func ReadAllS(src *os.File) (string, error) {
	bytes, err := ReadAll(src)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func FileCopy(filesystem fs.FS, src, dest string) error {
	var err error

	srcFile, err := filesystem.Open(src)
	defer srcFile.Close()
	if err != nil {
		return err
	}

	destFile, err := os.Create(dest)
	defer destFile.Close()
	if err != nil {
		return err
	}

	buf := make([]byte, 2048)
	for bytesRead, _ := srcFile.Read(buf); bytesRead != 0; bytesRead, _ = srcFile.Read(buf) {
		destFile.Write(buf)
	}

	return err
}

func RecursiveFileCopy(filesystem fs.FS, src, dest, dir string) error {
	if _, err := filesystem.Open(dest); errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(dest, 0700)
		if err != nil {
			return err
		}
	}
	var walkDir string
	if dir != "" {
		walkDir = src + "/" + dir
	} else {
		walkDir = src
	}

	err := fs.WalkDir(filesystem, walkDir,
		func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				return err
			}
			var destPath string = dest + path[len(src):]
			fmt.Printf("%s -> %s\n", path, destPath)

			if d.Name() == walkDir {
				return nil
			} else if d.IsDir() {
				return RecursiveFileCopy(filesystem, src, dest, path)
			} else {
				return FileCopy(filesystem, path, destPath)
			}
		},
	)

	return err
}

func ProcessInnerText(nodes []HtmlNode) {
	for i := range nodes {
		if nodes[i].Value != "" {
			var innerTextNodes TextNodeSlice
			innerTextNodes = append(innerTextNodes, TextNode{
				TextType: textTypeText,
				Text:     nodes[i].Value,
			})

			innerTextNodes, _ = innerTextNodes.SplitAll()
			if len(innerTextNodes) > 1 || innerTextNodes[0].TextType != textTypeText {
				for _, itNode := range innerTextNodes {
					child, _ := itNode.ToHTMLNode()
					nodes[i].Children = append(nodes[i].Children, child)
				}

				nodes[i].Value = ""
			}
		}

		if len(nodes[i].Children) > 0 {
			ProcessInnerText(nodes[i].Children)
		}

	}
}

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

	blocks := ParseMDBlocks(srcTxt)
	blockNodes, err := BlocksToHTMLNodes(blocks)
	if err != nil {
		return err
	}

	var pageTitle string

	ProcessInnerText(blockNodes)

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

func main() {
	var err error

	src := "static/index.md"
	template := "template.html"
	dest := "public/index.html"

	err = os.RemoveAll(dest)
	err = GeneratePage(src, template, dest)
	if err != nil {
		println(err.Error())
	}
}
