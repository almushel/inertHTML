package main

import (
	//"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type InertFlags struct {
	NoClobber, Interactive, Recursive, Verbose bool
}

func ErrPrintf(format string, a ...any) (int, error) {
	return fmt.Fprintf(os.Stderr, "inertHTML: "+format, a...)
}

func ErrPrintln(format string) (int, error) {
	return fmt.Fprintf(os.Stderr, "inertHTML: "+format+"\n")
}

func main() {
	var err error
	var flags InertFlags
	var src, template, dest string

	flag.BoolVar(&flags.NoClobber, "n", false, "do not overwrite an existing file")
	flag.BoolVar(&flags.Interactive, "i", false, "prompt before overwrite")
	flag.BoolVar(&flags.Recursive, "r", false, "process directories and their contents recursively")
	flag.BoolVar(&flags.Verbose, "v", false, "explain what is being done")
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
	src = flag.Arg(0)

	srcInfo, err := os.Stat(src)
	if err != nil {
		ErrPrintln(err.Error())
		os.Exit(1)
	} else if !srcInfo.IsDir() && !strings.HasSuffix(src, ".md") {
		ErrPrintln("Invalid source. Directory or *.md file expected")
		os.Exit(1)
	}

	if dest == "" {
		if strings.HasSuffix(src, ".md") {
			dest = src[:len(src)-len("md")] + "html"
		} else {
			dest = src
		}
	} else {
		destEnd := dest[strings.LastIndex(dest, "/")+1:]
		destIsDir := !(strings.Count(destEnd, ".") > 0)

		if !srcInfo.IsDir() && destIsDir {
			e := strings.LastIndex(src, "/") + 1
			dest += "/" + src[e:len(src)-len("md")] + "html"
		} else if srcInfo.IsDir() && !destIsDir {
			ErrPrintf("Cannot write output from directory %s to file %s\n", src, dest)
			os.Exit(1)
		}
	}

	if template != "" {
		err = ValidateTemplateFile(template)
		if err != nil {
			ErrPrintln(err.Error())
			os.Exit(1)
		}
	}

	if srcInfo.IsDir() {
		err = GenerateDirectory(src, template, dest, flags)
	} else {
		err = GeneratePageEx(src, template, dest, flags)
	}

	if err != nil {
		ErrPrintln(err.Error())
		os.Exit(1)
	}
}
