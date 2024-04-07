package main

import (
	//"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func ErrPrintf(format string, a ...any) (int, error) {
	return fmt.Fprintf(os.Stderr, "inertHTML: "+format, a...)
}

func ErrPrintln(format string) (int, error) {
	return fmt.Fprintf(os.Stderr, "inertHTML: "+format+"\n")
}

func main() {
	var err error

	var src, template, dest string
	flag.StringVar(&dest, "o", "", "output file/directory")
	//flag.StringVar(&template, "t", "", "html template for parsed markdown")
	flag.Parse()

	// NOTE: Maybe default to working directory if no src specified?
	if flag.NArg() < 1 {
		ErrPrintf("Source file/directory required\n")
		os.Exit(1)
	}
	src = flag.Arg(0)

	srcInfo, err := os.Stat(src)
	if !srcInfo.IsDir() && !strings.HasSuffix(src, ".md") {
		ErrPrintln("Invalid source. Directory or *.md file expected")
		os.Exit(1)
	} else if err != nil {
		ErrPrintf(err.Error())
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

	/*
		err = os.RemoveAll(dest)
		if err != nil {
			ErrPrintln(err.Error())
		}
	*/

	err = GeneratePage(src, template, dest)

	if err != nil {
		ErrPrintln(err.Error())
	}
}
