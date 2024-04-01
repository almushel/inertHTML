package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func RecursiveCopy(filesystem fs.FS, src, dest, dir string) error {
	fmt.Printf("Copying %s to %s\n", src, dest)

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
			println(path)
			var destPath string = dest + path[len(src):]

			if d.Name() == walkDir {
				return nil
			} else if d.IsDir() {
				errResult := os.MkdirAll(destPath, 0700)
				if errResult != nil {
					return errResult
				}
				return RecursiveCopy(filesystem, src, dest, dir+"/"+d.Name())
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

	err = RecursiveCopy(os.DirFS(wd), "test", "public", "")
	if err != nil {
		println(err.Error())
	}
}
