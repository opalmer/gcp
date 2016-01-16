package files

// copy.go - Top level module responsible for copying filess

import (
	"config"
	"github.com/op/go-logging"
	"github.com/ryanuber/go-glob"
	"os"
	"path/filepath"
)

var log = logging.MustGetLogger("gcp")

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

	if stat.Size() == 0 {
		return nil
	}

	exclude := false
	for _, exclusionPattern := range config.Exclude {
		if glob.Glob(exclusionPattern, path) {
			exclude = true
			break
		}
	}

	// See if there are any include statements that will force
	// the path to be included.
	if exclude {
		for _, inclusionPattern := range config.Include {
			if glob.Glob(inclusionPattern, path) {
				exclude = false
				break
			}
		}
	}

	if exclude {
		log.Debugf("Excluding %s", path)
		if stat.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	if !stat.IsDir() && stat.Size() > 0 {
		processing.Add(1)
		processPaths <- path
	}

	return nil
}

// Copy - The main function which does some preconfiguration and then passes
// off work to fileutil.Walk.
func Copy() {
	log.Infof("Copy %s -> %s", config.Source, config.Destination)

	processPaths = make(chan string)

	for worker := 1; worker <= config.Concurrency; worker++ {
		go process(processPaths)
	}

	if filepath.Walk(config.Source, Walk) != nil {
		log.Warningf("One or more failures while walking %s", config.Source)
	}

	processing.Wait()
}
