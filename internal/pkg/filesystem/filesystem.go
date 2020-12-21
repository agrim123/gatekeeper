package filesystem

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println("Filesystem error", err)
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return false
	case mode.IsRegular():
		return true
	}

	return false
}

func DoesExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func CopyFilesToDir(files []string, dst string) error {
	CreateDir(dst)
	for _, file := range files {
		Copy(file, dst+filepath.Base(file))
	}

	return nil
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	defer out.Close()

	err = out.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return err
	}

	return nil
}

func CreateDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			log.Fatal(errDir)
		}
	}
}
