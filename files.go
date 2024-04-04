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

func FileCopyRecursive(filesystem fs.FS, src, dest, dir string) error {
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
				return FileCopyRecursive(filesystem, src, dest, path)
			} else {
				return FileCopy(filesystem, path, destPath)
			}
		},
	)

	return err
}
