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
		fmt.Println(err)
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
	return out.Close()
}

func MoveFileToDir(path string) string {
	errDir := os.MkdirAll("/tmp/gatekeeper/keys", 0755)
	if errDir != nil {
		log.Fatal(errDir)
	}

	Copy(path, "/tmp/gatekeeper/keys/"+filepath.Base(path))
	return "/tmp/gatekeeper/keys"
}
