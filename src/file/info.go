package files

// info.go - Produces or constructs information for a given path.

import (
	"../config"
	"os"
	"path/filepath"
	"strings"
)

// Exists - Returns true if the given path exists
func Exists(path string) bool {
	_, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		return false
	}
	if err != nil {
		log.Fatal(err)
	}
	return true
}

// AbsolutePath - Returns the absolute path to the file.
func AbsolutePath(path string) string {
	absolute, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		log.Fatalf("Failed to determine absolute path for %s", path)
	}
	return absolute
}

// IsRelative - Returns true if the destination appears to be a
// subdirectory of the parent path.
func IsRelative(parent string, child string) bool {
	parentAbsolute := AbsolutePath(parent)
	childAbsolute := AbsolutePath(child)

	if parentAbsolute == childAbsolute {
		return true
	}

	relative, err := filepath.Rel(parentAbsolute, childAbsolute)
	if err != nil {
		log.Fatalf(
			"Failed to calculate relative paths for %s and %s: %s",
			parentAbsolute, childAbsolute, err)
	}

	// If the relative path starts with .. then Rel() had to walk
	// up and out of parentAbsolute.
	return !strings.HasPrefix(relative, "..")
}

// DestinationPath - Returns the path which the provided source should end
// up being copied to (minus the filename)
func DestinationPath(path string) string {
	if config.Compress {
		path += ".lzma"
	}

	if config.Encrypt {
		path += ".aes"
	}

	return filepath.Join(config.Destination, path)
}
