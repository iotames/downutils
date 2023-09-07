package service

import (
	"os"
	"time"
)

func GetFileModTime(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Now(), err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return time.Now(), err
	}
	return fi.ModTime(), nil
}
