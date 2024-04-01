package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func RecursiveCopy(filesystem fs.FS, src, dest, dir string) error {

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
				return RecursiveCopy(filesystem, src, dest, path)
			} else {
				srcFile, errResult := filesystem.Open(path)
				defer srcFile.Close()
				if errResult != nil {
					return errResult
				}

				destFile, errResult := os.Create(destPath)
				defer destFile.Close()
				if errResult != nil {
					return errResult
				}

				buf := make([]byte, 2048)
				for bytesRead, _ := srcFile.Read(buf); bytesRead != 0; bytesRead, _ = srcFile.Read(buf) {
					destFile.Write(buf)
				}

				return nil
			}
		},
	)

	return err
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	src := "static"
	dest := "public"

	err = os.RemoveAll(dest)
	err = RecursiveCopy(os.DirFS(wd), src, dest, "")
	if err != nil {
		println(err.Error())
	}
}
