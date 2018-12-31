package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

const workers = 32

type filepair struct {
	cardfile FileHash
	diridx   int
}

type stringpair struct {
	filename string
	reason   string
}

var (
	htmllog  os.File
	htmllock sync.Mutex
	treelock sync.RWMutex
	htmllist []stringpair
)

var done chan bool

func comparetrees(disktree Disktree, cardtree Cardtree) {
	log.Println("Comparing card and disk")

	htmllog, _ := os.Create("uncopied.html")
	htmllog.WriteString("<html>\n")
	htmllog.WriteString("<h1>files without backup on disk</h1>\n")
	htmllist = make([]stringpair, 0, 1000)

	channel := make(chan filepair, 32)
	done = make(chan bool)

	// create worker threads
	for i := 1; i < workers; i++ {
		go worker(i, channel, &disktree, &cardtree)
	}

	for cardpair, card := range cardtree.files {
		if diskimages, ok := disktree.files[cardpair]; ok {
			for idx, di := range diskimages {
				// check if picture on card is newer than on disk (we accept 1 hour difference due to summer/wintertime)
				if card.modtime.Sub(di.modtime) > time.Hour {
					fmt.Println(" no backup (based on date)", card.path, card.modtime)
					if verbose {
						fmt.Println("  candidate was", di.path, di.modtime)
					}
					htmllock.Lock()
					htmllist = append(htmllist, stringpair{card.path, "(based on date)"})
					htmllock.Unlock()
					if copydir != "" {
						fmt.Println("  copying", card.path, "to", filepath.Join(copydir, filepath.Base(card.path)))
					}
				} else {
					channel <- filepair{cardpair, idx}
				}
			}
		} else {
			fmt.Println(" no backup (based on name)", card.path)
			htmllock.Lock()
			htmllist = append(htmllist, stringpair{card.path, "(based on name)"})
			htmllock.Unlock()
			if copydir != "" {
				fmt.Println("  copying", card.path, "to", filepath.Join(copydir, filepath.Base(card.path)))
				CopyFile(card.path, filepath.Join(copydir, filepath.Base(card.path)), true)
			}
		}
	}

	close(channel)

	// wait for all workers
	if debug {
		log.Println("Waiting for workers...")
	}
	for i := 1; i < workers; i++ {
		<-done
	}

	sort.Slice(htmllist, func(i, j int) bool {
		return filepath.Base(htmllist[i].filename) < filepath.Base(htmllist[j].filename)
	})

	for _, f := range htmllist {
		htmllog.WriteString("<a href=\"" + "file:///" + f.filename + "\">" + f.filename + "</a> " + f.reason + "<br>\n")
	}

	htmllog.WriteString("</html>\n")
	htmllog.Close()

	log.Println("Comparison ended. See <uncopied.html> for list of files without copy.")
}

// worker to be run in parallel, to get parallel IO and sha256
func worker(me int, files chan filepair, disktree *Disktree, cardtree *Cardtree) {
	for f := range files {
		if debug {
			log.Println("... working on", f.cardfile, f.diridx)
		}

		treelock.RLock()
		di := (*disktree).files[f.cardfile][f.diridx]
		card := (*cardtree).files[f.cardfile]
		if !di.hashed {
			path := di.path
			treelock.RUnlock()
			hash := getsha256(path)
			treelock.Lock()
			(*disktree).files[f.cardfile][f.diridx].sha256 = hash
			(*disktree).files[f.cardfile][f.diridx].hashed = true
			treelock.Unlock()
			di = (*disktree).files[f.cardfile][f.diridx]
		} else {
			treelock.RUnlock()
			if debug {
				treelock.RLock()
				fmt.Println("... reusing hash", di.path)
				treelock.RUnlock()
			}
		}

		treelock.RLock()
		path := card.path
		treelock.RUnlock()

		hash := getsha256(path)

		treelock.Lock()
		card.sha256 = hash
		treelock.Unlock()

		treelock.RLock()
		if !bytes.Equal(card.sha256, di.sha256) {
			if debug {
				fmt.Println(" no backup (based on content)", card.path, di.path, me)
			} else {
				fmt.Println(" no backup (based on content)", card.path, di.path)
			}
			htmllock.Lock()
			htmllist = append(htmllist, stringpair{card.path, "(based on date)"})
			htmllock.Unlock()
			if copydir != "" {
				fmt.Println("  copying", card.path, "to", filepath.Join(copydir, filepath.Base(card.path)))
			}
		} else {
			if debug {
				fmt.Println("  same file", card.path, di.path, me)
			}
		}
		treelock.RUnlock()

	}

	if debug {
		log.Println("End of worker", me, ".")
	}
	done <- true
}

func getsha256(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("getsha256 1", err)
	}
	defer f.Close()

	h := sha256.New()
	stat, _ := f.Stat()
	if _, err := io.CopyN(h, f, min(stat.Size(), 1024*1024)); err != nil {
		log.Fatal("getsha256 2", err)
	}

	return h.Sum(nil)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
