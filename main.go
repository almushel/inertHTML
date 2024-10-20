package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/almushel/inertHTML/generator"
)

func ErrPrintf(format string, a ...any) (int, error) {
	return fmt.Fprintf(os.Stderr, "inertHTML: "+format, a...)
}

func ErrPrintln(format string) (int, error) {
	return fmt.Fprintf(os.Stderr, "inertHTML: "+format+"\n")
}

func main() {
	var err error
	var flags generator.InertFlags
	var src, template, dest string

	flag.BoolVar(&flags.NoClobber, "n", false, "do not overwrite an existing file")
	flag.BoolVar(&flags.Interactive, "i", false, "prompt before overwrite")
	flag.BoolVar(&flags.Recursive, "r", false, "process directories and their contents recursively")
	flag.BoolVar(&flags.Verbose, "v", false, "explain what is being done")
	flag.BoolVar(&flags.PagesAsDirs, "pagesAsDirs", false, "convert non-index source files to dest/index.html")
	flag.StringVar(&dest, "o", "", "write output to file/directory")
	flag.StringVar(&template, "t", "", "html template for parsed markdown")
	flag.Parse()

	if flags.NoClobber {
		flags.Interactive = false
	}

	if flag.NArg() < 1 {
		ErrPrintf("Source file/directory required\n")
		flag.Usage()
		os.Exit(1)
	}
	src = filepath.Clean(flag.Arg(0))

	srcInfo, err := os.Stat(src)
	if err != nil {
		ErrPrintln(err.Error())
		os.Exit(1)
	} else if !srcInfo.IsDir() && filepath.Ext(src) != ".md" {
		ErrPrintln("Invalid source. Directory or *.md file expected")
		os.Exit(1)
	}

	if dest == "" {
		if srcInfo.IsDir() {
			dest = src
		} else {
			dest = src[:len(src)-len("md")] + "html"
		}
	} else {
		var destIsDir bool
		dest := filepath.Clean(dest)

		destInfo, err := os.Stat(dest)
		if err == nil {
			destIsDir = destInfo.IsDir()
		} else if errors.Is(err, os.ErrNotExist) {
			// Assume path for nonexistent dest is either directory or .html (the only valid filetype)
			destIsDir = (filepath.Ext(dest) != ".html")
		} else {
			ErrPrintf(err.Error())
			os.Exit(1)
		}

		if !srcInfo.IsDir() && destIsDir {
			filename := filepath.Base(src)
			filename = filename[:len(filename)-len("md")] + "html"
			dest = filepath.Join(dest, filename)
		} else if srcInfo.IsDir() && !destIsDir {
			ErrPrintf("Cannot write output from directory %s to file %s\n", src, dest)
			os.Exit(1)
		}
	}

	if template != "" {
		err = generator.ValidateTemplateFile(template)
		if err != nil {
			ErrPrintln(err.Error())
			os.Exit(1)
		}
	}

	if srcInfo.IsDir() {
		err = generator.GenerateDirectory(src, template, dest, flags)
	} else {
		err = generator.GeneratePageEx(src, template, dest, flags)
	}

	if err != nil {
		ErrPrintln(err.Error())
		os.Exit(1)
	}
}
