package utils

import (
	"crypto/md5"
	"fmt"
	"time"
)

func GetCustomPath(path string) string {
	data := []byte(fmt.Sprintf("%s_%d", path, time.Now().Unix()))
	return fmt.Sprintf("%x", md5.Sum(data))
}
