package main

/*

	iscardcopied

	tool to verify if all pictures from a storage card are copied to a folder on disk

	A file does not have a copy on disktree
	- if name of file on card does not exist on disk (with same file size)
	- if name of file on card does exist on disk (with same file size), but date is newer on card
	- if name of file on card exists on disk, but they differ in contents

*/

import (
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

// list of file extensions to take care of
var extensionlist = []string{".TIF", ".TIFF", ".PSD", ".JPEG", ".ARW", ".MRW", ".JPG", ".CR", ".CR2", ".CR3", ".NEF",
	".MOV", ".MPEG", ".MPG", ".AVI", ".MP4", ".MKV", ".MTS", ".3GP"}

var extensions map[string]bool

var (
	verbose bool
	debug   bool
	copydir string
)

var opts struct {
	Card    string `long:"card" description:"card to read, including drive and path"`
	Disk    string `long:"disk" description:"disk folder to compare to"`
	Copy    string `long:"copy" description:"folder to copy new files to"`
	Verbose bool   `long:"verbose" short:"v" description:"verbose output, with doubles on disk"`
	Debug   bool   `long:"debug" short:"d" description:"debug output"`
}

func help() {
	fmt.Println("usage: iscardcopied --disk=<drive:\\folder> --card=<drive:>")
	fmt.Println(" checks if all relevant files (jpeg and raw and video) from card are somewhere on disk")
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}

	if opts.Disk == "" || opts.Card == "" {
		help()
		os.Exit(1)
	}

	if opts.Verbose {
		verbose = true
	} else {
		verbose = false
	}

	if opts.Debug {
		debug = true
	} else {
		debug = false
	}

	if opts.Copy != "" {
		copydir = opts.Copy
		stat, err := os.Stat(copydir)
		if os.IsNotExist(err) {
			fmt.Println("Copy target <" + copydir + "> does not exist!")
			os.Exit(2)
		}
		if !stat.IsDir() {
			fmt.Println("Copy target <" + copydir + "> is not a directory!")
			os.Exit(2)
		}
	}

	// create map with lower and uppercase extentions we care about
	extensions = make(map[string]bool)
	for _, e := range extensionlist {
		extensions[e] = true
		extensions[strings.ToLower(e)] = true
	}

	disktree := NewDisktree(opts.Disk)
	cardtree := NewCardtree(opts.Card)

	comparetrees(disktree, cardtree)
}
