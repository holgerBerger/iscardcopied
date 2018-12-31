package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Disktree maps filename to list of pathes
type Disktree struct {
	files map[FileHash][]FileEntry
}

// FileHash identifies a file with name and size
type FileHash struct {
	name string
	size int64
}

// FileEntry identifies a file with content
type FileEntry struct {
	path    string
	modtime time.Time
	hashed  bool
	sha256  []byte
}

// NewDisktree reads a directory tree on disk, creating a map from names to pathes
func NewDisktree(root string) Disktree {
	var disktree Disktree

	log.Println("Reading disk ", root)

	disktree.files = make(map[FileHash][]FileEntry)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if _, ok := extensions[ext]; ok {
				base := filepath.Base(path)
				pair := FileHash{base, info.Size()}
				if _, ok := disktree.files[pair]; ok {
					disktree.files[pair] = append(disktree.files[pair], FileEntry{path, info.ModTime(), false, []byte{0}})
					if verbose {
						fmt.Println("Same files?")
						for _, d := range disktree.files[pair] {
							fmt.Println("  ", d.path)
						}
					}
				} else {
					disktree.files[pair] = make([]FileEntry, 0, 2)
					disktree.files[pair] = append(disktree.files[pair], FileEntry{path, info.ModTime(), false, []byte{0}})
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path")
	}
	log.Println("Done Reading disk. Read", len(disktree.files), "files.")
	return disktree
}
