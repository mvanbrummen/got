package util

import "os"

func InitGotDir(gotDir string) error {
	return os.MkdirAll(gotDir, os.ModePerm)
}
