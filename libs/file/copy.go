package files

import (
	"github.com/op/go-logging"
	"os"
	"path/filepath"
)

var log = logging.MustGetLogger("gcp")

// Walker -
type Walker struct {
	src string
	dst string
}

var walker Walker

// SkipPath - Return True if the requested path should be skipped
func (w *Walker) SkipPath(path string) bool {
	return false
}

// Walk - Called by filepath.Walk to process files and directories
func Walk(path string, info os.FileInfo, err error) error {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return filepath.SkipDir
	} else if err != nil {
		log.Fatalf("Stat failed on %s (err: %s)", path, err)
	} else if stat.IsDir() || walker.SkipPath(path) {
		return nil
	}

	log.Debugf("Processing %s", path)

	return nil
}

// Copy - The main function which handles copy/compressing/encrypting files.
func Copy(
	dryRun bool, encryption bool, compression bool,
	source string, destination string) {

	log.Debug("Copy %s -> %s", source, destination)
	walker = Walker{source, destination}

	err := filepath.Walk(source, Walk)
	if err != nil {
		log.Warningf("One or more failures: %s", err)
	}

}
