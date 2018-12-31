// +build !windows

package main

import (
	"fmt"
	"io"
	"os"
)

// CopyFile generic implementation
func CopyFile(src, dst string, failoneexist bool) (int64, error) {
	srcfile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcfile.Close()

	srcfilestat, err := srcfile.Stat()
	if err != nil {
		return 0, err
	}

	if !srcfilestat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	dstfile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstfile.Close()
	return io.Copy(dstfile, srcfile)
}
