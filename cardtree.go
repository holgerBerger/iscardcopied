package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Cardtree maps filename to list of pathes
type Cardtree struct {
	files map[FileHash]FileEntry
}

// NewCardtree reads a directory tree on disk, creating a map from names to pathes
func NewCardtree(root string) Cardtree {
	var cardtree Cardtree

	log.Println("Reading card ", root)

	cardtree.files = make(map[FileHash]FileEntry)

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
				if _, ok := cardtree.files[pair]; ok {
					/// THIS SHOULD NOT HAPPEN ON CARD!!!!
					// cardtree.files[pair] = FileEntry{path, false, [32]byte{0}}
					fmt.Println("Same files on card?!?!!??")
				} else {
					cardtree.files[pair] = FileEntry{path, info.ModTime(), false, []byte{0}}
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path")
	}
	log.Println("Done Reading card. Read", len(cardtree.files), "files.")
	return cardtree
}
