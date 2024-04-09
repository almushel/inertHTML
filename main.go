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
	var flags InertFlags
	var src, template, dest string

	//flag.BoolVar(&flags.ForceCopy, "-f", false, "do not prompt before overwriting")
	flag.BoolVar(&flags.NoClobber, "n", false, "do not overwrite an existing file")
	flag.BoolVar(&flags.Interactive, "i", false, "prompt before overwrite")
	flag.StringVar(&dest, "o", "", "write output to file/directory")
	//flag.StringVar(&template, "t", "", "html template for parsed markdown")
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

	err = GeneratePageRecursiveEx(src, template, dest, flags)

	if err != nil {
		ErrPrintln(err.Error())
	}
}
