package files

import (
	"../config"
	"github.com/op/go-logging"
	"os"
	"path/filepath"
	"strings"
)

var log = logging.MustGetLogger("gcp")

const skipPath = 2
const keepPath = 1
const notMatched = 0

func skip(name string) int {
	for _, inclusion := range config.Include {
		if name == inclusion {
			return keepPath
		}

		matched, err := filepath.Match(inclusion, name)

		if err != nil {
			log.Fatalf(
				"filepath.Match('%s', '%s') failed (err: %b)",
				inclusion, name, err)
		}
		if matched {
			log.Debugf("Match(%s, %s)", name, inclusion)
			return keepPath
		}
	}

	for _, exclusion := range config.Exclude {
		if name == exclusion {
			return skipPath
		}

		matched, err := filepath.Match(exclusion, name)

		if err != nil {
			log.Fatalf(
				"filepath.Match('%s', '%s') failed (err: %s)",
				exclusion, name, err)
		}
		if matched {
			log.Debugf("Match(%s, %s)", name, exclusion)
			return skipPath
		}
	}

	return notMatched
}

// SkipPath - Return True if the requested path should be skipped
func SkipPath(path string) bool {
	result := skip(path)
	switch result {
	case skipPath:
		return true
	case keepPath:
		return false
	}

	// filepath.Match will **only** match non-Separator characters.  Because
	// of this, we split the path and try matching individual parts.
	for _, subpath := range strings.Split(path, string(os.PathSeparator)) {
		result := skip(subpath)
		switch result {
		case skipPath:
			return true
		case keepPath:
			return false
		case notMatched:
			continue
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
