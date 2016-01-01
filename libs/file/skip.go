package files

// skip.go - Responsible for determining if a given path should be skipped
// when copying files.

import (
	"../config"
	"os"
	"path/filepath"
	"strings"
)

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
