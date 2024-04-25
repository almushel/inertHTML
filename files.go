package main

import (
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

func ReadFileS(name string) (string, error) {
	buf, err := os.ReadFile(name)
	return string(buf), err
}

func FileCopy(src, dest string) error {
	var err error

	srcFile, err := os.Open(src)
	defer srcFile.Close()
	if err != nil {
		return err
	}

	destFile, err := CreateAll(dest)
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

func FileCopyRecursive(src, dest string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = fs.WalkDir(os.DirFS(wd), src,
		func(path string, d fs.DirEntry, err error) error {
			if d == nil || d.IsDir() {
				return err
			}

			var destPath string = dest + path[len(src):]

			return FileCopy(path, destPath)
		},
	)

	return err
}
