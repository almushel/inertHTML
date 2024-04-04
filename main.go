package main

import (
	"os"
)

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
