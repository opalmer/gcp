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
var wait sync.WaitGroup

// ProcessFile - Processes an individual file (called by Walk()).
func ProcessFile(path string) {
	destination := DestinationPath(path)
	log.Debug("%s -> %s", path, destination)
	output := NewOutput(path, destination)
	output.Process()
	log.Debug("[done] %s -> %s", path, destination)
	defer wait.Done()
}

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

	if !stat.IsDir() {
		wait.Add(1)
		go ProcessFile(path)
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
	wait.Wait()

}
