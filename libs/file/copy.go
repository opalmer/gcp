package files

import (
	"../config"
	"github.com/op/go-logging"
	"os"
	"path/filepath"
)

var log = logging.MustGetLogger("gcp")

// SkipPath - Return True if the requested path should be skipped
func SkipPath(path string) bool {
	// filepath.Match will **only** match non-Separator characters.  Because
	// of this, we split the path and try matching individual parts.
	basename := filepath.Base(path)
	dirname := filepath.Dir(path)

	for _, inclusionPattern := range config.Include {
		matched, err := filepath.Match(inclusionPattern, basename)
		if err == nil && matched || inclusionPattern == path ||
			inclusionPattern == dirname {
			return false
		}
	}

	for _, exclusionPattern := range config.Exclude {
		matched, err := filepath.Match(exclusionPattern, basename)
		if err == nil && matched || exclusionPattern == path ||
			exclusionPattern == dirname {
			return true
		}
	}

	return false
}

// DestinationPath - Returns the path which the provided source should end
// up being copied to (minus the filename)
func DestinationPath(path string) string {
	if config.Compress {
		path += ".7z"
	}

	if config.Encrypt {
		path += "aes"
	}

	return filepath.Join(config.Destination)
}

// ProcessFile - Processes an individual file (called by Walk()).
func ProcessFile(path string) {
	destination := DestinationPath(path)
	log.Debug("Processing %s -> %s", path, destination)

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

	ProcessFile(path)
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

}
