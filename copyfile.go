// +build windows

package main

import (
	"syscall"
	"unsafe"
)

var (
	modkernel32   = syscall.NewLazyDLL("kernel32.dll")
	procCopyFileW = modkernel32.NewProc("CopyFileW")
)

// CopyFile wraps windows function CopyFileW
func CopyFile(src, dst string, failIfExists bool) error {
	lpExistingFileName, err := syscall.UTF16PtrFromString(src)
	if err != nil {
		return err
	}

	lpNewFileName, err := syscall.UTF16PtrFromString(dst)
	if err != nil {
		return err
	}

	var bFailIfExists uint32
	if failIfExists {
		bFailIfExists = 1
	} else {
		bFailIfExists = 0
	}

	r1, _, err := syscall.Syscall(
		procCopyFileW.Addr(),
		3,
		uintptr(unsafe.Pointer(lpExistingFileName)),
		uintptr(unsafe.Pointer(lpNewFileName)),
		uintptr(bFailIfExists))

	if r1 == 0 {
		return err
	}
	return nil
}
