package files

import (
	"../config"
	"github.com/op/go-logging"
	"os"
	"path/filepath"
)

var log = logging.MustGetLogger("gcp")
var include []string
var exclude []string

// SkipPath - Return True if the requested path should be skipped
func SkipPath(path string) bool {
	// filepath.Match will **only** match non-Separator characters.  Because
	// of this, we split the path and try matching individual parts.
	basename := filepath.Base(path)
	dirname := filepath.Dir(path)

	for _, inclusionPattern := range include {
		matched, err := filepath.Match(inclusionPattern, basename)
		if err == nil && matched || inclusionPattern == path || inclusionPattern == dirname {
			return false
		}
	}

	for _, exclusionPattern := range exclude {
		matched, err := filepath.Match(exclusionPattern, basename)
		if err == nil && matched || exclusionPattern == path ||
			exclusionPattern == dirname {
			return true
		}
	}

	return false
}

// Walk - Called by filepath.Walk to process files and directories
func Walk(path string, info os.FileInfo, err error) error {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return filepath.SkipDir
	} else if err != nil {
		log.Fatalf("Stat failed on %s (err: %s)", path, err)
	} else if stat.IsDir() || SkipPath(path) {
		return nil
	}

	if !SkipPath(path) {
		log.Debug("Processing %s", path)
	}

	return nil
}

// Copy - The main function which handles copy/compressing/encrypting files.
func Copy(
	dryRun bool, encryption bool, compression bool,
	source string, destination string) {

	include = config.GetSlice("include")
	exclude = config.GetSlice("exclude")

	log.Debug("Copy %s -> %s", source, destination)

	err := filepath.Walk(source, Walk)
	if err != nil {
		log.Warningf("One or more failures: %s", err)
	}

}
