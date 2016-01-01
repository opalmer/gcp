package files

// copy.go - Top level module responsible for copying filess

import (
	"../config"
	"github.com/op/go-logging"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var log = logging.MustGetLogger("gcp")
var wait sync.WaitGroup
var processing = 0

// MaxThreads the maxiumum number of threads for handling files
var MaxThreads int

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

	if SkipPath(path) {
		if stat.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	for {
		if processing <= MaxThreads {
			break
		}
		time.Sleep(1)
	}

	if !stat.IsDir() {
		processing++
		destination := DestinationPath(path)
		output, err := NewOutput(path, destination)

		if err != nil {
			log.Fatalf("Failed to create output for %s", path)
		}
		wait.Add(1)
		go output.Process()
	}

	if processing >= MaxThreads {
		wait.Wait()
	}

	return nil
}

// Copy - The main function which does some preconfiguration and then passes
// off work to fileutil.Walk.
func Copy() {
	log.Infof("Copy %s -> %s", config.Source, config.Destination)

	err := filepath.Walk(config.Source, Walk)
	if err != nil {
		log.Warningf("One or more failures: %s", err)
	}

	// Make sure we wait on any remaining work
	defer wait.Wait()
}
