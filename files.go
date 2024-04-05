package main

import (
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

type WalkDirProc func(path string) error

func WalkDirRecursive(root, dir string, proc WalkDirProc) error {
	err := fs.WalkDir(os.DirFS(root), dir,
		func(path string, d fs.DirEntry, err error) error {
			if d == nil {
				return err
			}

			if path == dir {
				return err
			} else if d.IsDir() {
				return WalkDirRecursive(root, path, proc)
			} else {
				return proc(path)
			}
		},
	)

	return err
}

func FileCopyRecursive(src, dest string) error {
	proc := func(path string) error {
		var destPath string = dest + path[len(src):]

		return FileCopy(path, destPath)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return WalkDirRecursive(wd, src, proc)
}
