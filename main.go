// maildir-fix fixes maildirs missing some of the tmp, new and cur
// subdirectories.
//
// git removes directories that become empty, and maildir-using
// applications like mutt doesn't like that.
//
// Note: this is racy, but we can't do any better, with current git.
// Consider running maildir-fix after git pull, not concurrently with
// it.
//
// With the flag -depot=PATH, maildir-fix processes a whole directory
// of maildirs. Otherwise, list each maildir to fix on the command
// line.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s MAILDIR..\n", prog)
	fmt.Fprintf(os.Stderr, "  %s -depot=PATH\n", prog)
	flag.PrintDefaults()
}

var depot string

func init() {
	flag.StringVar(&depot, "depot", "", "mail depot containing maildirs to fix")
}

func processMaildir(logger *log.Logger, maildir string) (ok bool) {
	ok = true
	for _, child := range []string{"cur", "new", "tmp"} {
		err := os.Mkdir(filepath.Join(maildir, child), 0700)
		if err != nil && !os.IsExist(err) {
			logger.Printf("error processing maildir: %v", err)
			ok = false
		}
	}
	return
}

func processDepot(logger *log.Logger, depot string) (ok bool) {
	ok = true
	d, err := os.Open(depot)
	if err != nil {
		logger.Printf("error opening depot: %v", err)
		ok = false
		return
	}
	defer d.Close()

	for {
		fis, err := d.Readdir(1000)
		for _, fi := range fis {
			// all non-hidden subdirs of a depot are assumed to be
			// maildirs
			if fi.Name()[0] != '.' && fi.IsDir() {
				path := filepath.Join(depot, fi.Name())
				ok2 := processMaildir(logger, path)
				ok = ok && ok2
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Printf("error listing depot: %v", err)
			ok = false
			break
		}
	}
	return
}

func process(logger *log.Logger, depot string, maildirs []string) (ok bool) {
	ok = true
	if depot != "" {
		ok2 := processDepot(logger, depot)
		ok = ok && ok2
	}
	for _, maildir := range maildirs {
		ok2 := processMaildir(logger, maildir)
		ok = ok && ok2
	}
	return
}

func main() {
	logger := log.New(os.Stderr, prog+": ", 0)

	flag.Usage = usage
	flag.Parse()

	maildirs := flag.Args()
	if len(maildirs) == 0 && depot == "" {
		flag.Usage()
		os.Exit(2)
	}

	ok := process(logger, depot, maildirs)
	if !ok {
		os.Exit(1)
	}
}
