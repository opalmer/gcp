package files

// copy.go - Top level module responsible for copying filess

import (
	"../config"
	"github.com/op/go-logging"
	"os"
	"path/filepath"
	"sync"
)

var log = logging.MustGetLogger("gcp")
var filesProcessing = sync.WaitGroup{}

// Walk - Called by filepath.Walk to process files and directories.  This also
// performs the work of filtering paths which wish to process using SkipPath()
func Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Fatalf("Cannot process %s (err: %s)", path, err)
	}

	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return filepath.SkipDir
	} else if err != nil {
		log.Fatalf("Stat failed on %s (err: %s)", path, err)
	}

	if SkipPath(path) || stat.Size() == 0 {
		if stat.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	if !stat.IsDir() {
		filesProcessing.Add(1)
		channels.paths <- path
	}
	return nil
}

// Copy - The main function which does some preconfiguration and then passes
// off work to fileutil.Walk.
func Copy() {
	log.Infof("Copy %s -> %s", config.Source, config.Destination)

	// Before we walk over the paths, setup all the various
	// channels we'll use for processing.
	channels = Channels{
		paths: make(chan string),
		files: make(chan File)}

	for worker := 1; worker <= config.Concurrency; worker++ {
		go openpaths(channels.paths)
		go processfiles(channels.files)
	}

	if filepath.Walk(config.Source, Walk) != nil {
		log.Warningf("One or more failures while walking %s", config.Source)
	}

	filesProcessing.Wait()
}
