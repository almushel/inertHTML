package main

import (
	"os"
)

func main() {
	var err error

	src := "static"
	template := "template.html"
	dest := "public"

	err = os.RemoveAll(dest)
	err = GeneratePageRecursive(src, template, dest)
	if err != nil {
		println(err.Error())
	}
}
