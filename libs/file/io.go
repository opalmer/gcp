package files

// io.go - Responsible for performing input and output operations.

import (
	"../config"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const maxReadSize = 5e+7 // 50MB

// File - The main object used for storing and processing a single file.
type File struct {
	source   *os.File
	tempout  *os.File
	readsize int64
	outpath  string
}

// NewFile - Creates and returns File
func NewFile(sourcePath string) File {
	// Open the source file
	source, err := os.Open(sourcePath)
	if err != nil {
		if os.IsPermission(err) || os.IsNotExist(err) {
			log.Warningf("Could not open %s (err: %s)", sourcePath, err)
		} else {
			log.Fatalf("Failed to open %s (err: %s)", sourcePath, err)
		}
	}

	// Open the temporary output file.
	tempout, err := ioutil.TempFile(os.TempDir(), "gcp")
	if err != nil {
		source.Close()
		log.Fatal("Failed to open temporoary output file (err: %s)", err)
	}

	file := File{source: source, tempout: tempout}

	// Figure out how large of a slice we're supposed to
	// take when reading data.
	sourceStat, err := file.source.Stat()
	sourceSize := sourceStat.Size()

	if sourceSize < maxReadSize {
		file.readsize = sourceSize
	} else {
		file.readsize = maxReadSize
	}

	// Figure out what the final file path should be.
	outpath := filepath.Join(config.Destination, file.source.Name())

	if file.ShouldCompress() {
		outpath += ".lzma"
	}

	if file.ShouldEncrypt() {
		outpath += ".aes"
	}

	file.outpath = outpath

	return file
}

// ShouldCompress - Returns True if the file should be compressed.  Some kinds
// of files we shouldn't compress because it's either time consuming or we
// wouldn't gamin much by enabling compression.
func (file *File) ShouldCompress() bool {
	if !config.Compress {
		return false
	}

	name := strings.ToLower(file.source.Name())
	if strings.HasSuffix(name, ".iso") {
		return false
	}

	return true
}

// ShouldEncrypt - Returns True if the file should be compressed.  Some kinds
// of files we shouldn't compress because it's either time consuming or we
// wouldn't gamin much by enabling compression.
func (file *File) ShouldEncrypt() bool {
	if !config.Encrypt {
		return false
	}

	name := strings.ToLower(file.source.Name())
	if strings.HasSuffix(name, ".iso") {
		return false
	}

	return true
}

// DestinationExists - Returns True if the output path exists.
func (file *File) DestinationExists() bool {
	_, err := os.Stat(file.outpath)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else if os.IsPermission(err) {
		log.Fatalf("Failed to stat %s due to permission error", file.outpath)
	}
	log.Fatalf("Unhandled DestinationExists() for %s", file.outpath)
	return false
}

// Rename - Renames the temporary file
func (file *File) Rename() {
	log.Infof("%s -> %s", file.source.Name(), file.outpath)

	// Check if the parent directory exists, if not we'll
	// need to create it.
	err := os.MkdirAll(filepath.Dir(file.outpath), 0700)
	if err != nil {
		log.Fatal(
			"Failed to create parent directory for %s (err: %s)",
			file.outpath, err)
	}

	// If we're not performing compression or encryption then it's
	// a direct copy rather than a rename.
	if !file.ShouldCompress() && !file.ShouldEncrypt() {

	} else {
		err := os.Rename(file.tempout.Name(), file.outpath)
		if err != nil {
			log.Fatalf(
				"Failed to rename %s -> %s (err: %s)",
				file.tempout.Name(), file.outpath, err)
		}
	}

	_, err = file.tempout.Stat()
	if err == nil {
		err := os.Remove(file.tempout.Name())
		if err != nil {
			log.Warningf("Failed to remove %s", file.tempout.Name())
		}
	}
}

// Close - Closes the underlying files objects
func (file *File) Close() {
	errors := 0
	err := file.tempout.Sync()

	if err != nil {
		errors++
		log.Warningf("Failed to flush %s (err: %s)", file.tempout.Name(), err)
	}

	err = file.tempout.Close()
	if err != nil {
		errors++
		log.Warningf("Failed to close %s (err: %s)", file.tempout.Name(), err)
	}

	err = file.source.Close()
	if err != nil {
		errors++
		log.Warningf("Failed to close %s (err: %s)", file.source.Name(), err)
	}

	if errors > 0 {
		log.Fatal("One or more errors while calling Close()")
	}
}

// Done - Called when we've finished processing the requested file.
func (file *File) Done() {
	defer filesProcessing.Done()
	file.Close()
	file.Rename()
}
